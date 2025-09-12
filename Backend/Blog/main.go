package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/handler"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/repo"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default values.")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB.")

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "blog_db"
	}
	blogCollectionName  := os.Getenv("MONGO_COLLECTION_NAME")
	if blogCollectionName  == "" {
		blogCollectionName  = "blog"
	}

	blogCollection  := client.Database(dbName).Collection(blogCollectionName)

	blogRepo := repo.NewBlogRepository(blogCollection )
	blogService := service.NewBlogService(blogRepo)
	blogHandler := handler.NewBlogHandler(blogService)


	commentCollectionName := os.Getenv("COMMENT_COLLECTION_NAME")
	if commentCollectionName == "" {
		commentCollectionName = "comments"
	}

	commentCollection := client.Database(dbName).Collection(commentCollectionName)
	commentRepo := repo.NewCommentRepository(commentCollection)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService)



	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBlogServiceServer(grpcServer, blogHandler)
	pb.RegisterCommentServiceServer(grpcServer, commentHandler)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

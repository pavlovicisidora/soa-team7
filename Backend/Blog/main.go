package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/handler"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/repo"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/service"
	"github.com/pavlovicisidora/soa-team7/Backend/saga"
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
	blogCollectionName := os.Getenv("MONGO_COLLECTION_NAME")
	if blogCollectionName == "" {
		blogCollectionName = "blog"
	}

	natsURI := os.Getenv("NATS_URI")
	if natsURI == "" {
		natsURI = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURI)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	log.Println("Successfully connected to NATS.")

	blogCollection := client.Database(dbName).Collection(blogCollectionName)

	blogRepo := repo.NewBlogRepository(blogCollection)
	blogService := service.NewBlogService(blogRepo, nc)
	blogHandler := handler.NewBlogHandler(blogService)

	_, err = nc.Subscribe(saga.UserBlockedSubject, func(m *nats.Msg) {
		log.Printf("Received event on subject: %s", m.Subject)
		var event saga.UserBlockedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			return
		}
		if err := blogService.HandleUserBlocked(context.Background(), event.UserID); err != nil {
			log.Printf("SAGA failed in Blog service for user %s: %v. Publishing compensation event...", event.UserID, err)

			compensationEvent := saga.UserBlockFailedEvent{
				UserID: event.UserID,
				Reason: err.Error(),
			}
			eventData, _ := json.Marshal(compensationEvent)
			if pubErr := nc.Publish(saga.UserBlockFailedSubject, eventData); pubErr != nil {
				log.Printf("CRITICAL: Failed to publish compensation event: %v", pubErr)
			}
		} else {
			log.Printf("SAGA step successful in Blog service for user %s", event.UserID)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to NATS subject: %v", err)
	}

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

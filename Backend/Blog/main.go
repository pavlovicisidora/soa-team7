package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/handler"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/repo"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	collectionName := os.Getenv("MONGO_COLLECTION_NAME")
	if collectionName == "" {
		collectionName = "blog"
	}

	collection := client.Database(dbName).Collection(collectionName)

	blogRepo := repo.NewBlogRepository(collection)
	blogService := service.NewBlogService(blogRepo)
	blogHandler := handler.NewBlogHandler(blogService)

	router := mux.NewRouter()

	router.HandleFunc("/blog", blogHandler.CreateBlog).Methods("POST")
	router.HandleFunc("/blog", blogHandler.GetAllBlogs).Methods("GET")
	router.HandleFunc("/blog/{id}", blogHandler.GetBlogByID).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("Blog service listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

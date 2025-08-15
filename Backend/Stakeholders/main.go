package main

import (
	"context"
	"log" 
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pavlovicisidora/soa-team7/handler"
	"github.com/pavlovicisidora/soa-team7/repo"
	"github.com/pavlovicisidora/soa-team7/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func startServer(userHandler *handler.UserHandler, router *mux.Router) {

	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	router.HandleFunc("/userCreate", userHandler.Create).Methods("POST")


	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Println("Server starting on port :8081...")
	log.Fatal(http.ListenAndServe(":8081", corsObj(router)))

}

func main() {

	router := mux.NewRouter().StrictSlash(true)

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
		dbName = "stakeholders_db"
	}
	collectionName := os.Getenv("MONGO_COLLECTION_NAME")
	if collectionName == "" {
		collectionName = "stakeholders"
	}

	repo := repo.NewUserRepository(client, dbName, collectionName)
	service := &service.UserService{UserRepositroy: repo}
	handler := &handler.UserHandler{UserService: service}

	startServer(handler, router)

}

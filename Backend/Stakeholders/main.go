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

func startServer(userHandler *handler.UserHandler, profileHandler *handler.ProfileHandler, router *mux.Router) {

	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	router.HandleFunc("/register", userHandler.Create).Methods("POST")

	router.HandleFunc("/login", userHandler.Login).Methods("GET")

	router.HandleFunc("/blockUser/{username}", userHandler.BlockUser).Methods("PUT")

	router.HandleFunc("/usersInfo/", userHandler.FindAllInfo).Methods("GET")
	router.HandleFunc("/profiles/{userId}", profileHandler.FindByUserId).Methods("GET")
	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Println("Server starting on port :8081...")
	log.Fatal(http.ListenAndServe(":8081", corsObj(router)))

}

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
		dbName = "stakeholders_db"
	}
	collectionName := os.Getenv("MONGO_COLLECTION_NAME")
	if collectionName == "" {
		collectionName = "stakeholders"
	}

	userRepo := repo.NewUserRepository(client, dbName, collectionName)
	userService := &service.UserService{UserRepository: userRepo}
	userHandler := &handler.UserHandler{UserService: userService}

	profileService := &service.ProfileService{UserRepo: userRepo}
	profileHandler := &handler.ProfileHandler{ProfileService: profileService}
	router := mux.NewRouter().StrictSlash(true)

	startServer(userHandler, profileHandler, router)
}

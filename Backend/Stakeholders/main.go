package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/handler"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/repo"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
		dbName = "stakeholders_db"
	}
	collectionName := os.Getenv("MONGO_COLLECTION_NAME")
	if collectionName == "" {
		collectionName = "stakeholders"
	}

	userRepo := repo.NewUserRepository(client, dbName, collectionName)
	userService := &service.UserService{UserRepository: userRepo}
	profileService := &service.ProfileService{UserRepo: userRepo}

	// --- Pokretanje gRPC servera ---
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "8081"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Koristi novi gRPC handler
	grpcStakeholderServer := handler.NewStakeholderGRPCServer(*userService, *profileService)
	pb.RegisterStakeholderServiceServer(grpcServer, grpcStakeholderServer)

	reflection.Register(grpcServer)

	log.Printf("Stakeholders gRPC server starting on port :%s...", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

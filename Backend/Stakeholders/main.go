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
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/handler"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/repo"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/service"
	"github.com/pavlovicisidora/soa-team7/Backend/common/tracing"
	"github.com/pavlovicisidora/soa-team7/Backend/saga"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	otlpEndpoint := "jaeger:4317"

	tracerCloser, err := tracing.InitTracer("stakeholders-service", otlpEndpoint)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer tracerCloser.Close()

	err = godotenv.Load()
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

	userRepo := repo.NewUserRepository(client, dbName, collectionName)
	userService := &service.UserService{UserRepository: userRepo,
		NatsConn: nc}
	profileService := &service.ProfileService{UserRepo: userRepo}

	_, err = nc.Subscribe(saga.UserBlockFailedSubject, func(m *nats.Msg) {
		log.Printf("Received compensation event on subject: %s", m.Subject)
		var event saga.UserBlockFailedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("Error unmarshalling compensation event: %v", err)
			return
		}

		if err := userService.HandleBlockUserCompensation(context.Background(), event.UserID); err != nil {
			log.Printf("Failed to handle block user compensation for user %s: %v", event.UserID, err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to NATS subject: %v", err)
	}

	// --- Pokretanje gRPC servera ---
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "8081"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Koristi novi gRPC handler
	grpcStakeholderServer := handler.NewStakeholderGRPCServer(*userService, *profileService)
	pb.RegisterStakeholderServiceServer(grpcServer, grpcStakeholderServer)

	reflection.Register(grpcServer)

	log.Printf("Stakeholders gRPC server starting on port :%s...", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

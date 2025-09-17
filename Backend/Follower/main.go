package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/handler"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Follower/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/repo"
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/service"
	"github.com/pavlovicisidora/soa-team7/Backend/common/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	otlpEndpoint := "jaeger:4317"

	tracerCloser, err := tracing.InitTracer("follower-service", otlpEndpoint) // Prosleđujemo novi endpoint
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer tracerCloser.Close()
	// 1. Učitavanje .env fajla za konfiguraciju
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default values.")
	}

	// 2. Konekcija na Neo4j (sva logika je ovde)
	neo4jUri := os.Getenv("NEO4J_URI")
	if neo4jUri == "" {
		neo4jUri = "neo4j://localhost:7687"
	}
	neo4jUser := os.Getenv("NEO4J_USER")
	if neo4jUser == "" {
		neo4jUser = "neo4j" // Standardni username
	}
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")
	if neo4jPassword == "" {
		neo4jPassword = "follower" // Lozinka koju ste postavili u Neo4j Desktop-u
	}

	// Kreiranje Neo4j drajvera
	driver, err := neo4j.NewDriverWithContext(neo4jUri, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""))
	if err != nil {
		log.Fatalf("FATAL: Could not create Neo4j driver: %v", err)
	}
	defer driver.Close(context.Background())

	// Provera konekcije sa bazom
	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to Neo4j. Is the database running? Error: %v", err)
	}
	log.Println("Successfully connected to Neo4j.")

	// 3. Inicijalizacija slojeva (povezujemo sve delove)
	followerRepo := repo.NewNeo4jFollowRepository(driver)
	followerService := service.NewFollowService(followerRepo)
	followerHandler := handler.NewFollowerHandler(followerService)

	// 4. Podešavanje i pokretanje gRPC servera
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Koristite drugi port u odnosu na Blog servis
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	pb.RegisterFollowerServiceServer(grpcServer, followerHandler)

	log.Printf("Follower gRPC server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	tour_proto "github.com/pavlovicisidora/soa-team7/Backend/APIGateway/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Shopping/handler"
	"github.com/pavlovicisidora/soa-team7/Backend/Shopping/repo"
	"github.com/pavlovicisidora/soa-team7/Backend/Shopping/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	tourServiceAddress := os.Getenv("TOUR_SERVICE_ADDRESS")
	if tourServiceAddress == "" {
		tourServiceAddress = "localhost:9090"
	}
	conn, err := grpc.NewClient(tourServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	tourClient := tour_proto.NewTourGrpcServiceClient(conn)

	shoppingDb := db.Database("shopping_db")
	repo := repo.NewShoppingCartRepository(shoppingDb)
	service := service.NewShoppingCartService(repo, tourClient)
	handler := handler.NewShoppingCartHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/cart/checkout", handler.Checkout).Methods("POST")
	router.HandleFunc("/cart/{tourId}", handler.AddToCart).Methods("POST")
	router.HandleFunc("/cart/{tourId}", handler.RemoveFromCartHandler).Methods("DELETE")
	router.HandleFunc("/cart", handler.GetCartHandler).Methods("GET")
	router.HandleFunc("/token/{tourId}", handler.CheckTokenHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	log.Println("Shopping service starting on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/handler"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	tour_proto "github.com/pavlovicisidora/soa-team7/Backend/APIGateway/proto"
	blog_proto "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	blogServiceAddress := os.Getenv("BLOG_SERVICE_ADDRESS")
	if blogServiceAddress == "" {
		blogServiceAddress = "localhost:8082"
	}
	stakeholdersServiceAddress := os.Getenv("STAKEHOLDERS_SERVICE_ADDRESS")
	if stakeholdersServiceAddress == "" {
		stakeholdersServiceAddress = "localhost:8081"
	}

	tourServiceAddress := os.Getenv("TOUR_SERVICE_ADDRESS")
	if tourServiceAddress == "" {
		tourServiceAddress = "localhost:9090"
	}

	conn, err := grpc.NewClient(blogServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to blog service: %v", err)
	}
	defer conn.Close()

	connTour, err := grpc.NewClient(tourServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to tour service: %v", err)
	}
	defer connTour.Close()

	blogClient := blog_proto.NewBlogServiceClient(conn)
	blogHandler := handler.NewBlogHandler(blogClient)

	commentClient := blog_proto.NewCommentServiceClient(conn)
	commentHandler := handler.NewCommentHandler(commentClient)

	tourClient := tour_proto.NewTourGrpcServiceClient(connTour)
	tourHandler := handler.NewTourHandler(tourClient)

	stakeholdersURL, err := url.Parse("http://" + stakeholdersServiceAddress)
	if err != nil {
		log.Fatalf("Failed to parse stakeholders service URL: %v", err)
	}
	stakeholdersProxy := httputil.NewSingleHostReverseProxy(stakeholdersURL)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.Use(middleware.AuthMiddleware)

	apiRouter.PathPrefix("/blogs").Handler(http.StripPrefix("/api", blogHandler))
	apiRouter.PathPrefix("/comments").Handler(http.StripPrefix("/api", commentHandler))
	apiRouter.PathPrefix("/stakeholders").Handler(http.StripPrefix("/api", stakeholdersProxy))
	apiRouter.PathPrefix("/tours").Handler(http.StripPrefix("/api", tourHandler))
	// tourRouter := apiRouter.PathPrefix("/tours").Subrouter()
	// tourRouter.HandleFunc("/create", tourHandler.CreateTourHandler).Methods("POST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

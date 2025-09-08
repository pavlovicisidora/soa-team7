package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/handler"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	tour_proto "github.com/pavlovicisidora/soa-team7/Backend/APIGateway/proto"
	blog_proto "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
	follower_proto "github.com/pavlovicisidora/soa-team7/Backend/Follower/proto"
	stakeholders_proto "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
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

	followerServiceAddress := os.Getenv("FOLLOWER_SERVICE_ADDRESS")
	if followerServiceAddress == "" {
		followerServiceAddress = "localhost:8084"
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

<<<<<<< HEAD
	connStakeholders, err := grpc.NewClient(stakeholdersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to tour service: %v", err)
	}
	defer connStakeholders.Close()

	blogClient := blog_proto.NewBlogServiceClient(conn)
	blogHandler := handler.NewBlogHandler(blogClient)

	commentClient := blog_proto.NewCommentServiceClient(conn)
	commentHandler := handler.NewCommentHandler(commentClient)

=======
>>>>>>> 4b8e782692281132203f8047d5dee74ea330f633
	followerConn, err := grpc.NewClient(followerServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to follower service: %v", err)
	}
	defer followerConn.Close()

	followerClient := follower_proto.NewFollowerServiceClient(followerConn)
	followerHandler := handler.NewFollowerHandler(followerClient)

	blogClient := blog_proto.NewBlogServiceClient(conn)
	blogHandler := handler.NewBlogHandler(blogClient, followerClient)

	commentClient := blog_proto.NewCommentServiceClient(conn)
	commentHandler := handler.NewCommentHandler(commentClient)

	tourClient := tour_proto.NewTourGrpcServiceClient(connTour)
	tourHandler := handler.NewTourHandler(tourClient)

	keyPointClient := tour_proto.NewKeyPointGrpcServiceClient(connTour)
	keyPointHandler := handler.NewKeyPointHandler(keyPointClient)

	reviewClient := tour_proto.NewReviewGrpcServiceClient(connTour)
	reviewHandler := handler.NewReviewHandler(reviewClient)

	//stakeholdersURL, err := url.Parse("http://" + stakeholdersServiceAddress)
	//if err != nil {
	//	log.Fatalf("Failed to parse stakeholders service URL: %v", err)
	//}
	//stakeholdersProxy := httputil.NewSingleHostReverseProxy(stakeholdersURL)

	///NOVOO
	stakeholdersClient := stakeholders_proto.NewStakeholderServiceClient(connStakeholders)
	userHandler := handler.NewAPIUserHandler(stakeholdersClient)
	////

	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("./uploads/"))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))

	publicApiRouter := router.PathPrefix("/api").Subrouter()
	publicApiRouter.HandleFunc("/images/upload", handler.UploadImageHandler).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.Use(middleware.AuthMiddleware)

	apiRouter.PathPrefix("/blogs").Handler(http.StripPrefix("/api", blogHandler))
	apiRouter.PathPrefix("/comments").Handler(http.StripPrefix("/api", commentHandler))
	apiRouter.PathPrefix("/users").Handler(http.StripPrefix("/api", userHandler))

	apiRouter.PathPrefix("/follower").Handler(http.StripPrefix("/api", followerHandler))

	apiRouter.PathPrefix("/tours").Handler(http.StripPrefix("/api", tourHandler))
	apiRouter.PathPrefix("/keypoints").Handler(http.StripPrefix("/api", keyPointHandler))
	apiRouter.PathPrefix("/reviews").Handler(http.StripPrefix("/api", reviewHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

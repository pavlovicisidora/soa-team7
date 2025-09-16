package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/handlers"
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

	shoppingServiceAddress := os.Getenv("SHOPPING_SERVICE_ADDRESS")
	if shoppingServiceAddress == "" {
		shoppingServiceAddress = "localhost:8085"
	}
	shoppingURL, _ := url.Parse("http://" + shoppingServiceAddress)
	shoppingProxy := httputil.NewSingleHostReverseProxy(shoppingURL)
	shoppingProxy.Director = func(req *http.Request) {
		originalReq := req

		req.URL.Scheme = shoppingURL.Scheme
		req.URL.Host = shoppingURL.Host
		req.URL.Path = originalReq.URL.Path

		userID := originalReq.Context().Value(middleware.UserKey)
		if userID != nil {
			req.Header.Set("X-User-ID", userID.(string))
		}
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

	stakeholdersGrpcAddress := "stakeholders-server:8089"
	stakeholdersConn, err := grpc.NewClient(stakeholdersGrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to stakeholders gRPC service: %v", err)
	}
	defer stakeholdersConn.Close()

	connStakeholders, err := grpc.NewClient(stakeholdersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to tour service: %v", err)
	}
	defer connStakeholders.Close()

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

	//stakeholdersURL, err := url.Parse("http://" + stakeholdersServiceAddress)
	//if err != nil {
	//	log.Fatalf("Failed to parse stakeholders service URL: %v", err)
	//}
	//stakeholdersProxy := httputil.NewSingleHostReverseProxy(stakeholdersURL)

	///NOVOO
	stakeholdersClient := stakeholders_proto.NewStakeholderServiceClient(connStakeholders)
	userHandler := handler.NewAPIUserHandler(stakeholdersClient)
	profileHandler := handler.NewProfileHandler(stakeholdersClient)
	reviewHandler := handler.NewReviewHandler(reviewClient, stakeholdersClient)
	////

	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("./uploads/"))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))

	publicApiRouter := router.PathPrefix("/api").Subrouter()
	publicApiRouter.HandleFunc("/images/upload", handler.UploadImageHandler).Methods("POST")
	publicApiRouter.HandleFunc("/users/login", userHandler.LoginHandler).Methods("POST")
	publicApiRouter.HandleFunc("/users/register", userHandler.CreateUserHandler).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)

	apiRouter.PathPrefix("/blogs").Handler(http.StripPrefix("/api", blogHandler))
	apiRouter.PathPrefix("/comments").Handler(http.StripPrefix("/api", commentHandler))
	apiRouter.PathPrefix("/users").Handler(http.StripPrefix("/api", userHandler))
	apiRouter.PathPrefix("/follower").Handler(http.StripPrefix("/api", followerHandler))
	apiRouter.PathPrefix("/tours").Handler(http.StripPrefix("/api", tourHandler))
	apiRouter.PathPrefix("/keypoints").Handler(http.StripPrefix("/api", keyPointHandler))
	apiRouter.PathPrefix("/reviews").Handler(http.StripPrefix("/api", reviewHandler))
	apiRouter.PathPrefix("/profile").Handler(http.StripPrefix("/api", profileHandler))
	apiRouter.PathPrefix("/shopping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/shopping")
		shoppingProxy.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	corsRouter := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)

	log.Printf("API Gateway starting on port %s", port)
	if err := http.ListenAndServe(":"+port, corsRouter); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	follower_proto "github.com/pavlovicisidora/soa-team7/Backend/Follower/proto"
)


type FollowerHandler struct {
	client follower_proto.FollowerServiceClient
}


func NewFollowerHandler(client follower_proto.FollowerServiceClient) *FollowerHandler {
	return &FollowerHandler{client: client}
}


func (h *FollowerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/follower/follow/{userId}", h.FollowUserHandler).Methods("POST")
	router.HandleFunc("/follower/following", h.GetMyFollowingHandler).Methods("GET")
	router.HandleFunc("/follower/recommendations", h.GetRecommendationsHandler).Methods("GET")

	router.ServeHTTP(w, r)
}


func (h *FollowerHandler) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	
	followerID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || followerID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	
	vars := mux.Vars(r)
	followedID := vars["userId"]
	if followedID == "" {
		http.Error(w, "Followed user ID is required in the path", http.StatusBadRequest)
		return
	}

	
	grpcReq := &follower_proto.FollowUserRequest{
		FollowerId: followerID,
		FollowedId: followedID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()


	_, err := h.client.FollowUser(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to follow user via gRPC: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	
	w.WriteHeader(http.StatusNoContent)
}

func (h *FollowerHandler) GetMyFollowingHandler(w http.ResponseWriter, r *http.Request) {
	
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	
	grpcReq := &follower_proto.GetFollowingRequest{UserId: userID}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	
	resp, err := h.client.GetFollowing(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to get following list via gRPC: %v", err)
		
		http.Error(w, "Failed to retrieve following list", http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetUsers())
}

func (h *FollowerHandler) GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	
	grpcReq := &follower_proto.GetFollowRecommendationsRequest{UserId: userID}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	
	resp, err := h.client.GetFollowRecommendations(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to get follow recommendations via gRPC: %v", err)
		http.Error(w, "Failed to retrieve recommendations", http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetUsers())
}

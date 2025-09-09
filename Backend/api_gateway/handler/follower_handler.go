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

// FollowerHandler je HTTP handler koji komunicira sa Follower mikroservisom.
type FollowerHandler struct {
	client follower_proto.FollowerServiceClient
}

// NewFollowerHandler kreira novu instancu handlera.
func NewFollowerHandler(client follower_proto.FollowerServiceClient) *FollowerHandler {
	return &FollowerHandler{client: client}
}

// ServeHTTP registruje rute za follower funkcionalnosti.
func (h *FollowerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/follower/follow/{userId}", h.FollowUserHandler).Methods("POST")
	router.HandleFunc("/follower/following", h.GetMyFollowingHandler).Methods("GET")
	router.HandleFunc("/follower/recommendations", h.GetRecommendationsHandler).Methods("GET")

	router.ServeHTTP(w, r)
}

// FollowUserHandler obrađuje zahtev za praćenje drugog korisnika.
func (h *FollowerHandler) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Dobijamo ID ulogovanog korisnika iz konteksta (postavio ga je middleware)
	followerID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || followerID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Dobijamo ID korisnika kojeg treba pratiti iz URL-a
	vars := mux.Vars(r)
	followedID := vars["userId"]
	if followedID == "" {
		http.Error(w, "Followed user ID is required in the path", http.StatusBadRequest)
		return
	}

	// 3. Kreiramo gRPC zahtev
	grpcReq := &follower_proto.FollowUserRequest{
		FollowerId: followerID,
		FollowedId: followedID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 4. Pozivamo Follower mikroservis
	_, err := h.client.FollowUser(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to follow user via gRPC: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	// 5. Vraćamo uspešan odgovor bez tela
	w.WriteHeader(http.StatusNoContent)
}

func (h *FollowerHandler) GetMyFollowingHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Dobijamo ID ulogovanog korisnika iz konteksta
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Kreiramo gRPC zahtev sa dobijenim ID-jem
	grpcReq := &follower_proto.GetFollowingRequest{UserId: userID}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 3. Pozivamo GetFollowing RPC na Follower mikroservisu
	resp, err := h.client.GetFollowing(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to get following list via gRPC: %v", err)
		// U produkciji, ovde biste mogli mapirati gRPC greške (npr. codes.NotFound)
		// u odgovarajuće HTTP statuse (npr. http.StatusNotFound).
		http.Error(w, "Failed to retrieve following list", http.StatusInternalServerError)
		return
	}

	// 4. Vraćamo uspešan odgovor sa JSON telom
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetUsers())
}

func (h *FollowerHandler) GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Dobijamo ID ulogovanog korisnika iz konteksta
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// 2. Kreiramo gRPC zahtev
	grpcReq := &follower_proto.GetFollowRecommendationsRequest{UserId: userID}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 3. Pozivamo GetFollowRecommendations RPC na Follower mikroservisu
	resp, err := h.client.GetFollowRecommendations(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to get follow recommendations via gRPC: %v", err)
		http.Error(w, "Failed to retrieve recommendations", http.StatusInternalServerError)
		return
	}

	// 4. Vraćamo uspešan odgovor sa listom preporučenih korisnika
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetUsers())
}

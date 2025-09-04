package handler

import (
	"context"
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

	router.HandleFunc("/follow/{userId}", h.FollowUserHandler).Methods("POST")

	/*// Ruta za otpraćivanje korisnika: DELETE /users/{id}/follow
	router.HandleFunc("/users/{id}/follow", h.UnfollowUserHandler).Methods("DELETE")

	// Ruta za dobijanje preporuka za praćenje za ulogovanog korisnika: GET /users/recommendations
	router.HandleFunc("/users/recommendations", h.GetRecommendationsHandler).Methods("GET")

	// Bonus: Rute za dobijanje liste pratilaca i praćenih
	router.HandleFunc("/users/{id}/following", h.GetFollowingHandler).Methods("GET")
	router.HandleFunc("/users/{id}/followers", h.GetFollowersHandler).Methods("GET")*/

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

/*
// UnfollowUserHandler obrađuje zahtev za otpraćivanje korisnika.
func (h *FollowerHandler) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || followerID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followedID := vars["id"]
	if followedID == "" {
		http.Error(w, "Followed user ID is required in the path", http.StatusBadRequest)
		return
	}

	grpcReq := &follower_proto.UnfollowUserRequest{
		FollowerId: followerID,
		FollowedId: followedID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.client.UnfollowUser(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to unfollow user via gRPC: %v", err)
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetRecommendationsHandler vraća preporuke za praćenje za ulogovanog korisnika.
func (h *FollowerHandler) GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	grpcReq := &follower_proto.GetFollowRecommendationsRequest{
		UserId: userID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetFollowRecommendations(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to get recommendations via gRPC: %v", err)
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetUsers())
}


// GetFollowingHandler vraća listu korisnika koje dati korisnik prati.
func (h *FollowerHandler) GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
    userID := mux.Vars(r)["id"]
    if userID == "" {
        http.Error(w, "User ID is required in the path", http.StatusBadRequest)
        return
    }

    grpcReq := &follower_proto.GetFollowingRequest{UserId: userID}

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    resp, err := h.client.GetFollowing(ctx, grpcReq)
    if err != nil {
        log.Printf("Failed to get following list via gRPC: %v", err)
        http.Error(w, "Failed to get following list", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp.GetUsers())
}

// GetFollowersHandler vraća listu korisnika koji prate datog korisnika.
func (h *FollowerHandler) GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
    userID := mux.Vars(r)["id"]
    if userID == "" {
        http.Error(w, "User ID is required in the path", http.StatusBadRequest)
        return
    }

    grpcReq := &follower_proto.GetFollowersRequest{UserId: userID}

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    resp, err := h.client.GetFollowers(ctx, grpcReq)
    if err != nil {
        log.Printf("Failed to get followers list via gRPC: %v", err)
        http.Error(w, "Failed to get followers list", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp.GetUsers())
}*/

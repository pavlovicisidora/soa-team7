package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileHandler struct {
	client pb.StakeholderServiceClient
}

func NewProfileHandler(client pb.StakeholderServiceClient) *ProfileHandler {
	return &ProfileHandler{
		client: client,
	}
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	protected := router.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/profile", h.GetUserProfileById).Methods("GET")
	protected.HandleFunc("/profile", h.PatchProfile).Methods("PATCH")

	router.ServeHTTP(w, r)
}

func (h *ProfileHandler) GetUserProfileById(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userIDFromToken == "" {
		http.Error(w, "Could not retrieve user ID from token", http.StatusInternalServerError)
		return
	}

	grpcRequest := &pb.GetUserProfileByIdRequest{
		UserId: userIDFromToken,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetUserProfileById(ctx, grpcRequest)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			http.Error(w, "Profile not found for the logged-in user", http.StatusNotFound)
		} else {
			log.Printf("Error fetching current user profile: %v", err)
			http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *ProfileHandler) PatchProfile(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userIDFromToken == "" {

		http.Error(w, "Could not retrieve user ID from token", http.StatusInternalServerError)
		return
	}

	var grpcRequest pb.PatchProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&grpcRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcRequest.UserId = userIDFromToken

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.client.PatchProfile(ctx, &grpcRequest)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.InvalidArgument {
			http.Error(w, st.Message(), http.StatusBadRequest)
		} else {
			log.Printf("Error updating profile: %v", err)
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

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
)

// APIUserHandler je REST adapter za korisničke operacije
type APIUserHandler struct {
	client pb.StakeholderServiceClient
}

func NewAPIUserHandler(client pb.StakeholderServiceClient) *APIUserHandler {
	return &APIUserHandler{
		client: client,
	}
}

// ServeHTTP kreira rute za sve korisničke operacije
func (h *APIUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	// Public rute (bez JWT)
	router.HandleFunc("/users/public", h.GetUserPublicInfoHandler).Methods("GET")
	router.HandleFunc("/users/login", h.LoginHandler).Methods("POST")
	router.HandleFunc("/users/register", h.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/user", h.GetUser).Methods("GET")
	router.HandleFunc("/users/position", h.UpdateUserPosition).Methods("PUT")

	// Zaštićene rute sa JWT middleware
	protected := router.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/users", h.GetAllUsersHandler).Methods("GET")
	protected.HandleFunc("/users/block", h.BlockUserHandler).Methods("POST")

	router.ServeHTTP(w, r)
}

// GetUserPublicInfoHandler vraća javne informacije o korisniku
func (h *APIUserHandler) GetUserPublicInfoHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetUserPublicInfo(ctx, &pb.GetUserPublicInfoRequest{UserId: username})
	if err != nil {
		log.Printf("Error fetching public info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// LoginHandler omogućava login korisnika i vraća JWT token
// LoginHandler omogućava login korisnika
func (h *APIUserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Login(ctx, &req)
	if err != nil {
		log.Printf("Login failed: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateUserHandler registruje novog korisnika
func (h *APIUserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Create(ctx, &req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *APIUserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value("userRole").(string)
	if !ok || userRole == "" {
		http.Error(w, "Could not retrieve user role from context", http.StatusInternalServerError)
		return
	}

	if userRole != "ADMIN" {
		http.Error(w, "Forbidden: an admin role is required to access this resource", http.StatusForbidden)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetAllUsers(ctx, &pb.GetAllUsersRequest{})
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Users)
}

// BlockUserHandler blokira korisnika (zaštićeno JWT)
func (h *APIUserHandler) BlockUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.BlockUser(ctx, &pb.BlockUserRequest{Username: username})
	if err != nil {
		log.Printf("Error blocking user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.ProtoReflect())
}

func (h *APIUserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetUser(ctx, &pb.GetUserRequest{UserId: userID})
	if err != nil {
		http.Error(w, "Failed to get profile via gRPC", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetUser())
}

func (h *APIUserHandler) UpdateUserPosition(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(string)

	var reqBody struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.client.UpdateUserPosition(ctx, &pb.UpdateUserPositionRequest{
		UserId: userID,
		Lat:    reqBody.Lat,
		Long:   reqBody.Long,
	})
	if err != nil {
		http.Error(w, "Failed to update position via gRPC", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "position updated"})
}

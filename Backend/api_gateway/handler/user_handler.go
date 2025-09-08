package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
)

type APIUserHandler struct {
	GrpcClient pb.StakeholderServiceClient
}

func NewAPIUserHandler(client pb.StakeholderServiceClient) *APIUserHandler {
	return &APIUserHandler{GrpcClient: client}
}

// ServeHTTP kreira rute za sve korisničke operacije
func (h *APIUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/users/public", h.GetUserPublicInfoHandler).Methods("GET")
	router.HandleFunc("/users/login", h.LoginHandler).Methods("POST")
	router.HandleFunc("/users", h.GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/users", h.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/block", h.BlockUserHandler).Methods("POST")

	router.ServeHTTP(w, r)
}

// GetUserPublicInfoHandler vraća javne informacije o korisniku
func (h *APIUserHandler) GetUserPublicInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.GrpcClient.GetUserPublicInfo(ctx, &pb.GetUserPublicInfoRequest{UserId: userID})
	if err != nil {
		log.Printf("Error fetching public info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// LoginHandler omogućava login korisnika
func (h *APIUserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req pb.Login
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.GrpcClient.Login(ctx, &req)
	if err != nil {
		log.Printf("Login failed: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetAllUsersHandler vraća listu svih korisnika
func (h *APIUserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.GrpcClient.GetAllUsers(ctx, &pb.GetAllUsersRequest{})
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Users)
}

// CreateUserHandler kreira novog korisnika
func (h *APIUserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.GrpcClient.CreateUser(ctx, &req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BlockUserHandler blokira korisnika po username-u
func (h *APIUserHandler) BlockUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.GrpcClient.BlockUser(ctx, &pb.BlockUserRequest{Username: username})
	if err != nil {
		log.Printf("Error blocking user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

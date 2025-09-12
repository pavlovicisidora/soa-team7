package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	tour_proto "github.com/pavlovicisidora/soa-team7/Backend/APIGateway/proto"
)

type CreateTourRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Difficulty  string   `json:"difficulty"`
	Tags        []string `json:"tags"`
}

type TourHandler struct {
	client tour_proto.TourGrpcServiceClient
}

func NewTourHandler(client tour_proto.TourGrpcServiceClient) *TourHandler {
	return &TourHandler{client: client}
}

func (h *TourHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/tours", h.CreateTour).Methods("POST")
	router.HandleFunc("/tours", h.GetAllToursById).Methods("GET")
	router.HandleFunc("/tours/all", h.GetAllTours).Methods("GET")
	router.HandleFunc("/tours/{id}", h.GetTourById).Methods("GET")

	router.ServeHTTP(w, r)
}
func (h *TourHandler) CreateTour(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	role := r.Context().Value("userRole").(string)
	if role != "VODIC" {
		http.Error(w, "Forbidden: only VODIC can create tour.", http.StatusForbidden)
		return
	}

	var reqBody CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcRequest := &tour_proto.CreateTourRequest{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Difficulty:  reqBody.Difficulty,
		Tags:        reqBody.Tags,
		AuthorId:    userID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.CreateTour(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to create tour via gRPC: %v", err)
		http.Error(w, "Failed to create tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp.GetTour())
}
func (h *TourHandler) GetAllToursById(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserKey).(string)

	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	role := r.Context().Value("userRole").(string)
	if role != "VODIC" {
		http.Error(w, "Forbidden: only VODIC can see tours.", http.StatusForbidden)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := h.client.GetAllToursById(ctx, &tour_proto.GetAllToursByIdRequest{AuthorId: userID})
	if err != nil {
		log.Printf("Failed to get all tours: %v", err)
		http.Error(w, "Failed to get all tours", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetTours())
}
func (h *TourHandler) GetAllTours(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserKey).(string)

	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := h.client.GetAllTours(ctx, &tour_proto.GetAllToursRequest{})
	if err != nil {
		log.Printf("Failed to get all tours: %v", err)
		http.Error(w, "Failed to get all tours", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetTours())
}

func (h *TourHandler) GetTourById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Tour ID is missing in path", http.StatusBadRequest)
		return
	}
	tourID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetTourById(ctx, &tour_proto.GetTourByIdRequest{TourId: int32(tourID)})
	if err != nil {
		log.Printf("Failed to get tour by ID: %v", err)
		http.Error(w, "Failed to get tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetTour())
}

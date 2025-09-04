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

type CreateKeyPointRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageUrl    string  `json:"image_url"`
}

type KeyPointHandler struct {
	client tour_proto.KeyPointGrpcServiceClient
}

func NewKeyPointHandler(client tour_proto.KeyPointGrpcServiceClient) *KeyPointHandler {
	return &KeyPointHandler{client: client}
}

func (h *KeyPointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/keypoints/{id}", h.CreateKeyPoint).Methods("POST")
	router.HandleFunc("/keypoints/{id}", h.GetKeyPointsTour).Methods("GET")

	router.ServeHTTP(w, r)
}

func (h *KeyPointHandler) CreateKeyPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	role := r.Context().Value("userRole").(string)
	if role != "VODIC" {
		http.Error(w, "Forbidden: only VODIC can create keypoints.", http.StatusForbidden)
		return
	}
	var reqBody CreateKeyPointRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tourID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}
	grpcRequest := &tour_proto.CreateKeyPointRequest{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Latitude:    reqBody.Latitude,
		Longitude:   reqBody.Longitude,
		ImageUrl:    reqBody.ImageUrl,
		TourId:      int32(tourID),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.CreateKeyPoint(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to create tour via gRPC: %v", err)
		http.Error(w, "Failed to create tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp.GetKeypoint())
}

func (h *KeyPointHandler) GetKeyPointsTour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID, ok := r.Context().Value(middleware.UserKey).(string)

	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	role := r.Context().Value("userRole").(string)
	if role != "VODIC" {
		http.Error(w, "Forbidden: only VODIC can see keypoints.", http.StatusForbidden)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	tourID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}
	resp, err := h.client.GetKeyPointsForTour(ctx, &tour_proto.GetKeyPointsForTourRequest{TourId: int32(tourID)})
	if err != nil {
		log.Printf("Failed to get all tours: %v", err)
		http.Error(w, "Failed to get all tours", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetKeyPoints())
}

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
	router.HandleFunc("/tours", h.GetUpdateTour).Methods("PUT")

	router.HandleFunc("/tours/{id}/start", h.StartTour).Methods("POST")
	router.HandleFunc("/tours/execution/{execId}/abandon", h.AbandonTour).Methods("POST")
	router.HandleFunc("/tours/execution/{execId}/complete", h.CompleteTour).Methods("POST")
	router.HandleFunc("/tours/execution/{execId}", h.GetTourExecution).Methods("GET")

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

func (h *TourHandler) StartTour(w http.ResponseWriter, r *http.Request) {
	touristID := r.Context().Value(middleware.UserKey).(string)
	tourIDStr := mux.Vars(r)["id"]
	tourID, _ := strconv.ParseInt(tourIDStr, 10, 32)
	var reqBody struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	json.NewDecoder(r.Body).Decode(&reqBody)

	grpcRequest := &tour_proto.StartTourRequest{
		TourId:    int32(tourID),
		TouristId: touristID,
		Latitude:  reqBody.Latitude,
		Longitude: reqBody.Longitude,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.StartTour(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to start tour via gRPC: %v", err)
		http.Error(w, "Failed to start tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp.GetTourExecution())
}

func (h *TourHandler) AbandonTour(w http.ResponseWriter, r *http.Request) {
	execIDStr := mux.Vars(r)["execId"]
	execID, _ := strconv.ParseInt(execIDStr, 10, 32)

	grpcRequest := &tour_proto.AbandonTourRequest{TourExecutionId: int32(execID)}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.AbandonTour(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to abandon tour via gRPC: %v", err)
		http.Error(w, "Failed to abandon tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetTourExecution())
}

func (h *TourHandler) CompleteTour(w http.ResponseWriter, r *http.Request) {
	execIDStr := mux.Vars(r)["execId"]
	execID, _ := strconv.ParseInt(execIDStr, 10, 32)

	grpcRequest := &tour_proto.CompleteTourRequest{TourExecutionId: int32(execID)}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.CompleteTour(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to complete tour via gRPC: %v", err)
		http.Error(w, "Failed to complete tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetTourExecution())
}

func (h *TourHandler) GetTourExecution(w http.ResponseWriter, r *http.Request) {
	execIDStr := mux.Vars(r)["execId"]
	execID, _ := strconv.ParseInt(execIDStr, 10, 32)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetTourExecution(ctx, &tour_proto.GetTourExecutionRequest{TourExecutionId: int32(execID)})
	if err != nil {
		http.Error(w, "Failed to get tour execution", http.StatusInternalServerError)
		return
	}

	responsePayload := map[string]interface{}{
		"id":                resp.TourExecution.Id,
		"tour_id":           resp.TourExecution.TourId,
		"tourist_id":        resp.TourExecution.TouristId,
		"status":            resp.TourExecution.Status,
		"tour":              resp.Tour,
		"start_time":        resp.TourExecution.StartTime,
		"completition_time": resp.TourExecution.CompletionTime,
		"last_activity":     resp.TourExecution.LastActivity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func (h *TourHandler) GetUpdateTour(w http.ResponseWriter, r *http.Request) {
	
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	role := r.Context().Value("userRole").(string)
	if role != "VODIC" {
		http.Error(w, "Forbidden: only VODIC can update a tour.", http.StatusForbidden)
		return
	}

	
	var reqBody tour_proto.Tour
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	
	if reqBody.GetId() == 0 { 
		http.Error(w, "Tour ID must be provided in the request body", http.StatusBadRequest)
		return
	}
	
	
	if userID != reqBody.GetAuthorId() {
		http.Error(w, "Forbidden: You can only update your own tours.", http.StatusForbidden)
		return
	}


	
	grpcRequest := &tour_proto.UpdateTourRequest{
		Tour: &reqBody,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	
	resp, err := h.client.UpdateTour(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to update tour via gRPC: %v", err)
		http.Error(w, "Failed to update tour", http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetTour())
}

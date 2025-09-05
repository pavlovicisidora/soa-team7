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

type CreateReviewRequest struct {
	Rating       int32    `json:"rating"`
	Comment      string   `json:"comment"`
	VisitingDate string   `json:"visitingdate"`
	Images       []string `json:"images"`
}

type ReviewHandler struct {
	client tour_proto.ReviewGrpcServiceClient
}

func NewReviewHandler(client tour_proto.ReviewGrpcServiceClient) *ReviewHandler {
	return &ReviewHandler{client: client}
}
func (h *ReviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/reviews/{id}", h.CreateReview).Methods("POST")
	router.HandleFunc("/reviews/{id}", h.GetAllReviewsForTour).Methods("GET")

	router.ServeHTTP(w, r)
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	tourID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	var reqBody CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcRequest := &tour_proto.CreateReviewRequest{
		Rating:    reqBody.Rating,
		Comment:   reqBody.Comment,
		TouristId: "123",
		VisitDate: reqBody.VisitingDate,
		Images:    reqBody.Images,
		TourId:    int32(tourID),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.CreateReview(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to create review via gRPC: %v", err)
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp.GetReview())
}

func (h *ReviewHandler) GetAllReviewsForTour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	tourID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(middleware.UserKey).(string)

	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := h.client.GetReviewsForTour(ctx, &tour_proto.GetReviewForTourRequest{TourId: int32(tourID)})
	if err != nil {
		log.Printf("Failed to get all reviews: %v", err)
		http.Error(w, "Failed to get all reviews", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetReviews())
}

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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KeyPointRequest struct {
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

	router.HandleFunc("/keypoints/{tourId}", h.CreateKeyPoint).Methods("POST")
	router.HandleFunc("/keypoints/{tourId}", h.GetKeyPointsTour).Methods("GET")
	router.HandleFunc("/keypoints/{id}", h.UpdateKeyPoint).Methods("PUT")
	router.HandleFunc("/keypoints/{id}", h.DeleteKeyPoint).Methods("DELETE")

	router.ServeHTTP(w, r)
}

func (h *KeyPointHandler) CreateKeyPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["tourId"]
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

	var reqBody KeyPointRequest
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

		log.Printf("Failed to create keypoint via gRPC: %v", err)
		http.Error(w, "Failed to create keypoint", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp.GetKeypoint())
}

func (h *KeyPointHandler) GetKeyPointsTour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["tourId"]
	// userID, ok := r.Context().Value(middleware.UserKey).(string)

	// if !ok || userID == "" {
	// 	http.Error(w, "User ID not found in context", http.StatusUnauthorized)
	// 	return
	// }

	/*role := r.Context().Value("userRole").(string)
	if role != "VODIC" {
		http.Error(w, "Forbidden: only VODIC can see keypoints.", http.StatusForbidden)
		return
	}*/
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	tourID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}
	resp, err := h.client.GetKeyPointsForTour(ctx, &tour_proto.GetKeyPointsForTourRequest{TourId: int32(tourID)})
	if err != nil {

		log.Printf("Failed to get all keypoints: %v", err)
		http.Error(w, "Failed to get all keypoints", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetKeyPoints())
}

func (h *KeyPointHandler) UpdateKeyPoint(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	role := r.Context().Value("userRole").(string)
	if role != "VODIC" {
		http.Error(w, "Forbidden: only VODIC can update keypoints.", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	keyPointID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid keypoint ID", http.StatusBadRequest)
		return
	}

	var reqBody KeyPointRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcRequest := &tour_proto.UpdateKeyPointRequest{
		Id:          int32(keyPointID),
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Latitude:    reqBody.Latitude,
		Longitude:   reqBody.Longitude,
		ImageUrl:    reqBody.ImageUrl,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.UpdateKeyPoint(ctx, grpcRequest)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			http.Error(w, st.Message(), http.StatusNotFound)
		} else {
			log.Printf("Failed to update keypoint via gRPC: %v", err)
			http.Error(w, "Failed to update keypoint", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetKeypoint())
}
func (h *KeyPointHandler) DeleteKeyPoint(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	role := r.Context().Value("userRole").(string)
	if role != "VODIC" {
		http.Error(w, "Forbidden: only VODIC can delete keypoints.", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	keyPointID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid keypoint ID", http.StatusBadRequest)
		return
	}

	grpcRequest := &tour_proto.DeleteKeyPointRequest{Id: int32(keyPointID)}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err = h.client.DeleteKeyPoint(ctx, grpcRequest)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			http.Error(w, st.Message(), http.StatusNotFound)
		} else {
			log.Printf("Failed to delete keypoint via gRPC: %v", err)
			http.Error(w, "Failed to delete keypoint", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

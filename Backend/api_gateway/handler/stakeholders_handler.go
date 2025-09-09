package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	stakeholders_proto "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
)

func (h *StakeholderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/stakeholders/user", h.GetUser).Methods("GET")
	router.HandleFunc("/stakeholders/position", h.UpdateUserPosition).Methods("PUT")

	router.ServeHTTP(w, r)
}

type StakeholderHandler struct {
	client stakeholders_proto.StakeholderServiceClient
}

func NewStakeholderHandler(client stakeholders_proto.StakeholderServiceClient) *StakeholderHandler {
	return &StakeholderHandler{client: client}
}

func (h *StakeholderHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetUser(ctx, &stakeholders_proto.GetUserRequest{UserId: userID})
	if err != nil {
		http.Error(w, "Failed to get profile via gRPC", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.GetUser())
}

func (h *StakeholderHandler) UpdateUserPosition(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.client.UpdateUserPosition(ctx, &stakeholders_proto.UpdateUserPositionRequest{
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

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileHandler struct {
	ProfileService *service.ProfileService
}

func (handler *ProfileHandler) FindByUserId(writer http.ResponseWriter, req *http.Request) {
	userIdStr := mux.Vars(req)["userId"]
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {

		http.Error(writer, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	profile, err := handler.ProfileService.GetUserProfile(req.Context(), userId)
	if err != nil {
		http.Error(writer, "Profile not found", http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(profile)
}

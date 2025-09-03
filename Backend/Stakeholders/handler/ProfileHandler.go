package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/auth"
	"github.com/pavlovicisidora/soa-team7/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileHandler struct {
	ProfileService *service.ProfileService
}

func (handler *ProfileHandler) FindByUserId(writer http.ResponseWriter, req *http.Request) {
	tokenStr := req.Header.Get("Authorization") // "Bearer <token>"
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	claims, err := auth.VerifyJWT(tokenStr)
	if err != nil {
		http.Error(writer, "Unathorized"+err.Error(), http.StatusUnauthorized)
		return
	}

	if claims.Role != "TURISTA" && claims.Role != "VODIC" {
		http.Error(writer, "Forbidden: only TURISTA AND VODIC can see user profile", http.StatusForbidden)
		return
	}
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
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(profile)
}

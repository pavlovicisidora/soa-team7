package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/service"
	"github.com/pavlovicisidora/soa-team7/auth"
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

func (handler *ProfileHandler) PatchProfile(w http.ResponseWriter, r *http.Request) {
	// 1. JWT iz headera
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	claims, err := auth.VerifyJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// 2. Pretvori userID iz tokena u ObjectID
	userId, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 3. Parsiraj samo polja koja korisnik želi da menja
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 4. Pozovi servis za partial update
	if err := handler.ProfileService.UpdateUserProfileFields(r.Context(), userId, updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

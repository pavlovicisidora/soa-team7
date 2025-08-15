package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pavlovicisidora/soa-team7/model"
	"github.com/pavlovicisidora/soa-team7/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func (h *UserHandler) GetAllUsers(v http.ResponseWriter, r *http.Request) {
	log.Println("Test?")
}

func (handler *UserHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.UserService.Create(req.Context(), &user)
	if err != nil {
		http.Error(writer, "Error while creating a new user", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	
}

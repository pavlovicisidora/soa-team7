package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pavlovicisidora/soa-team7/model"
	"github.com/pavlovicisidora/soa-team7/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func (handler *UserHandler) GetAllUsers(writer http.ResponseWriter, req *http.Request) {
	users, err := handler.UserService.GetAllUsers(req.Context())
	if err != nil {
		http.Error(writer, "Error while collecting all users", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK) 

	// Encode users slice into JSON and write to response
	if err := json.NewEncoder(writer).Encode(users); err != nil {
		http.Error(writer, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
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
	println("Succesfully added user!")
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")

}

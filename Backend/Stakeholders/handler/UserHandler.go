package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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

	if user.Role != "VODIC" && user.Role != "TURISTA" {
		http.Error(writer, "Invalid role. Role must be VODIC or TURISTA.", http.StatusBadRequest)
		return
	}

	user.Blocked = false

	err = handler.UserService.Create(req.Context(), &user)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "mail") {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, "DB error", http.StatusInternalServerError)
		}
		return
	}

	println("Succesfully added user!")
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")

}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (handler *UserHandler) Login(writer http.ResponseWriter, req *http.Request) {
	var reqBody loginRequest
	err := json.NewDecoder(req.Body).Decode(&reqBody)

	if err != nil {
		http.Error(writer, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := handler.UserService.Login(req.Context(), reqBody.Username, reqBody.Password)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Blocked {
		http.Error(writer, "Your account has been blocked.", http.StatusForbidden)
		return
	}

	println("Succesfull login!")
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}

func (handler *UserHandler) BlockUser(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	err := handler.UserService.BlockUser(req.Context(), username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent) // 204 No Content
}
func (handler *UserHandler) FindAllInfo(writer http.ResponseWriter, req *http.Request) {
	users, err := handler.UserService.FindAllInfo(req.Context())
	if err != nil {
		http.Error(writer, "Error while collecting all users", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(users); err != nil {
		http.Error(writer, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

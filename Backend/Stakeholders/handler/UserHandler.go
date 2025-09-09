package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/auth"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/model"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedStakeholderServiceServer
	UserService    *service.UserService
	ProfileService *service.ProfileService
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

	token, err := auth.GenerateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		http.Error(writer, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"username": user.Username,
			"role":     user.Role,
		},
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)

	println("Succesfull login!")
}

func (handler *UserHandler) BlockUser(writer http.ResponseWriter, req *http.Request) {

	//Uzimamo token iz heder-a
	tokenStr := req.Header.Get("Authorization") // "Bearer <token>"
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	//Verifikujemo token, autentifikacija, da li je token uopste validan(ispravan potpis, nije istekao, ispravan format)
	claims, err := auth.VerifyJWT(tokenStr)
	if err != nil {
		http.Error(writer, "Unathorized"+err.Error(), http.StatusUnauthorized)
		return
	}

	//Proveravamo ulogu, autorizacija
	if claims.Role != "ADMIN" {
		http.Error(writer, "Forbidden: only ADMIN can block users", http.StatusForbidden)
		return
	}

	vars := mux.Vars(req)
	username := vars["username"]

	err = handler.UserService.BlockUser(req.Context(), username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (handler *UserHandler) FindAllInfo(writer http.ResponseWriter, req *http.Request) {
	tokenStr := req.Header.Get("Authorization") // "Bearer <token>"
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	claims, err := auth.VerifyJWT(tokenStr)
	if err != nil {
		http.Error(writer, "Unathorized"+err.Error(), http.StatusUnauthorized)
		return
	}

	if claims.Role != "ADMIN" {
		http.Error(writer, "Forbidden: only ADMIN can see users information", http.StatusForbidden)
		return
	}
	id := claims.UserID
	users, err := handler.UserService.FindAllInfo(req.Context(), id)
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
func (handler *UserHandler) GetUserPublicInfo(ctx context.Context, req *pb.GetUserPublicInfoRequest) (*pb.GetUserPublicInfoResponse, error) {
	log.Printf("Received GetUserPublicInfo request for user ID: %s", req.GetUserId())

	userIDStr := req.GetUserId()
	if userIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "User ID cannot be empty")
	}

	objectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		log.Printf("Invalid ObjectID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User ID format")
	}

	user, err := handler.UserService.FindById(ctx, objectID)

	if err != nil {
		log.Printf("User not found for ID %s: %v", userIDStr, err)
		return nil, status.Errorf(codes.NotFound, "User with id %s not found", userIDStr)
	}
	profile, err := handler.ProfileService.GetUserProfile(ctx, objectID)
	if err != nil {
		log.Printf("Profile not found for ID %s: %v", userIDStr, err)
		return nil, status.Errorf(codes.NotFound, "User with id %s not found", userIDStr)
	}

	response := &pb.GetUserPublicInfoResponse{
		UserId:        user.ID.Hex(),
		Username:      user.Username,
		Name:          profile.Name,
		Surname:       profile.Surname,
		ProfilePicUrl: profile.ProfilePic,
	}

	return response, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User ID format")
	}

	user, err := h.UserService.FindById(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	protoUser := &pb.User{
		Id:        user.ID.Hex(),
		Username:  user.Username,
		Mail:      user.Mail,
		Role:      user.Role,
		Latitude:  user.Latitude,
		Longitude: user.Longitude,
	}

	return &pb.GetUserResponse{User: protoUser}, nil
}

func (h *UserHandler) UpdateUserPosition(ctx context.Context, req *pb.UpdateUserPositionRequest) (*pb.UpdateUserPositionResponse, error) {
	err := h.UserService.UpdateUserPosition(ctx, req.GetUserId(), req.GetLat(), req.GetLong())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update position: %v", err)
	}
	return &pb.UpdateUserPositionResponse{Status: "position updated"}, nil
}

package handler

import (
	"context"
	"log"
	"strings"

	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/auth"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/model"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StakeholderGRPCServer implementira gRPC StakeholderServiceServer interfejs.
type StakeholderGRPCServer struct {
	pb.UnimplementedStakeholderServiceServer
	UserService    service.UserService
	ProfileService service.ProfileService
}

// NewStakeholderGRPCServer kreira novi StakeholderGRPCServer.
func NewStakeholderGRPCServer(userService service.UserService, profileService service.ProfileService) *StakeholderGRPCServer {
	return &StakeholderGRPCServer{
		UserService:    userService,
		ProfileService: profileService,
	}
}

// toProtoUser konvertuje model.User u pb.User
func toProtoUser(user *model.User) *pb.User {
	if user == nil {
		return nil
	}
	return &pb.User{
		UserId:   user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
		Blocked:  user.Blocked,
	}
}

// toProtoPublicUser konvertuje model.User i model.Profile u pb.GetUserPublicInfoResponse
func toProtoPublicUser(user *model.User, profile *model.Profile) *pb.GetUserPublicInfoResponse {
	if user == nil {
		return nil
	}
	response := &pb.GetUserPublicInfoResponse{
		UserId:   user.ID.Hex(),
		Username: user.Username,
		// Ostali podaci iz profila ako postoji
	}
	if profile != nil {
		response.Name = profile.Name
		response.Surname = profile.Surname
		response.ProfilePicUrl = profile.ProfilePic
	}
	return response
}

// Login implementira gRPC Login metodu.
func (s *StakeholderGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("gRPC Login request for username: %s", req.GetUsername())

	user, err := s.UserService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		// Konvertuj greške u gRPC status kodove
		if strings.Contains(err.Error(), "Invalid credentials") {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
		}
		if strings.Contains(err.Error(), "User not found") {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Login failed: %v", err)
	}

	// if user.Blocked {
	// 	return nil, status.Errorf(codes.PermissionDenied, "Account is blocked")
	// }

	token, err := auth.GenerateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token: %v", err)
	}

	log.Printf("gRPC Login successful for user: %s", user.Username)
	return &pb.LoginResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

// Create implementira gRPC Create metodu.
func (s *StakeholderGRPCServer) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("gRPC CreateUser request for username: %s", req.GetUsername())

	user := &model.User{
		Username: req.GetUsername(),
		Password: req.GetPassword(), // Lozinka bi trebalo da bude hešovana u servisu!
		Role:     req.GetRole(),
		Blocked:  false,
	}

	// Proveri validnost uloge pre kreiranja
	if user.Role != "VODIC" && user.Role != "TURISTA" && user.Role != "ADMIN" { // Dodao sam i ADMIN za svaki slučaj
		return nil, status.Errorf(codes.InvalidArgument, "Invalid role. Role must be VODIC, TURISTA or ADMIN.")
	}

	err := s.UserService.Create(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") { // Ili za email
			return nil, status.Errorf(codes.AlreadyExists, "User with username %s already exists", req.GetUsername())
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	log.Printf("gRPC CreateUser successful for user: %s", user.Username)
	return &pb.CreateUserResponse{
		UserId: user.ID.Hex(),
	}, nil
}

// GetAllUsers implementira gRPC GetAllUsers metodu.
func (s *StakeholderGRPCServer) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	log.Println("gRPC GetAllUsers request received.")

	users, err := s.UserService.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch all users: %v", err)
	}

	var protoUsers []*pb.User
	for _, user := range users {
		protoUsers = append(protoUsers, toProtoUser(&user))
	}

	return &pb.GetAllUsersResponse{Users: protoUsers}, nil
}

// GetUserPublicInfo implementira gRPC GetUserPublicInfo metodu.
func (s *StakeholderGRPCServer) GetUserPublicInfo(ctx context.Context, req *pb.GetUserPublicInfoRequest) (*pb.GetUserPublicInfoResponse, error) {
	log.Printf("gRPC GetUserPublicInfo request for user ID: %s", req.GetUserId())

	userIDStr := req.GetUserId()
	if userIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "User ID cannot be empty")
	}

	objectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		log.Printf("Invalid ObjectID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User ID format")
	}

	user, err := s.UserService.FindById(ctx, objectID)
	if err != nil {
		log.Printf("User not found for ID %s: %v", userIDStr, err)
		return nil, status.Errorf(codes.NotFound, "User with id %s not found", userIDStr)
	}

	profile, err := s.ProfileService.GetUserProfile(ctx, objectID)
	if err != nil {
		log.Printf("Profile not found for ID %s (or no public profile data): %v", userIDStr, err)
		// Nije kritična greška ako nema profila, samo vratimo korisnika bez profilnih podataka
		return toProtoPublicUser(user, nil), nil // Vrati samo korisnika ako nema profila
	}

	return toProtoPublicUser(user, profile), nil
}

// BlockUser implementira gRPC BlockUser metodu.
func (s *StakeholderGRPCServer) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.BlockUserResponse, error) {
	log.Printf("gRPC BlockUser request for username: %s", req.GetUsername())

	// Implementacija autorizacije za gRPC
	// U gRPC-u se autorizacija obično radi preko interceptora ili provere metadate u kontekstu
	// Za sada, pretpostavljamo da je poziv već autorizovan ako dolazi iz API Gateway-a
	// Ako želiš da proveriš ulogu, moraš je proslediti u metadati gRPC poziva.
	// Npr: metadata, ok := metadata.FromIncomingContext(ctx) ... i proveri header

	username := req.GetUsername()
	if username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Username cannot be empty")
	}

	err := s.UserService.BlockUser(ctx, username)
	if err != nil {
		if strings.Contains(err.Error(), "User not found") {
			return nil, status.Errorf(codes.NotFound, "User %s not found", username)
		}
		return nil, status.Errorf(codes.Internal, "Failed to block user %s: %v", username, err)
	}

	log.Printf("gRPC BlockUser successful for user: %s", username)
	return &pb.BlockUserResponse{}, nil
}

func (h *StakeholderGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User ID format")
	}

	user, err := h.UserService.FindById(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	protoUser := &pb.UserPS{
		Id:        user.ID.Hex(),
		Username:  user.Username,
		Mail:      user.Mail,
		Role:      user.Role,
		Latitude:  user.Latitude,
		Longitude: user.Longitude,
	}

	return &pb.GetUserResponse{User: protoUser}, nil
}

func (h *StakeholderGRPCServer) UpdateUserPosition(ctx context.Context, req *pb.UpdateUserPositionRequest) (*pb.UpdateUserPositionResponse, error) {
	err := h.UserService.UpdateUserPosition(ctx, req.GetUserId(), req.GetLat(), req.GetLong())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update position: %v", err)
	}
	return &pb.UpdateUserPositionResponse{Status: "position updated"}, nil
}

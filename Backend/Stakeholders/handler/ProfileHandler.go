package handler

import (
	"context"
	"log"

	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/model"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileHandler struct {
	ProfileService *service.ProfileService
	pb.UnimplementedStakeholderServiceServer
}

func toProtoProfile(profile *model.Profile) *pb.GetUserProfileByIdResponse {
	if profile == nil {
		return nil
	}
	return &pb.GetUserProfileByIdResponse{
		Name:          profile.Name,
		Surname:       profile.Surname,
		ProfilePicUrl: profile.ProfilePic,
		Bio:           profile.Bio,
		Motto:         profile.Motto,
	}
}

func (s *StakeholderGRPCServer) GetUserProfileById(ctx context.Context, req *pb.GetUserProfileByIdRequest) (*pb.GetUserProfileByIdResponse, error) {
	log.Printf("gRPC GetUserProfileById request for user ID: %s", req.GetUserId())

	userIdStr := req.GetUserId()
	if userIdStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "User ID cannot be empty")
	}
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID format")
	}

	profile, err := s.ProfileService.GetUserProfile(ctx, userId)
	if err != nil {
		log.Printf("Profile not found for ID %s: %v", userIdStr, err)
		return nil, status.Errorf(codes.NotFound, "Profile not found for user ID: %s", userIdStr)
	}
	log.Printf("Successfully fetched profile for user ID: %s", userIdStr)
	return toProtoProfile(profile), nil
}
func (s *StakeholderGRPCServer) PatchProfile(ctx context.Context, req *pb.PatchProfileRequest) (*pb.PatchProfileResponse, error) {
	log.Println("gRPC PatchProfile request received.")

	userId, err := primitive.ObjectIDFromHex(req.GetUserId())
	if err != nil {
		log.Printf("ERROR: Invalid user ID format in PatchProfile: %s", req.GetUserId())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID in token")
	}

	updates := make(map[string]interface{})

	if req.GetName() != "" {
		updates["name"] = req.GetName()
	}
	if req.GetSurname() != "" {
		updates["surname"] = req.GetSurname()
	}
	if req.GetProfilePicUrl() != "" {
		updates["picture"] = req.GetProfilePicUrl()
	}
	if req.GetBio() != "" {
		updates["bio"] = req.GetBio()
	}
	if req.GetMotto() != "" {
		updates["motto"] = req.GetMotto()
	}

	if len(updates) == 0 {
		log.Printf("WARN: PatchProfile request for user %s with no fields to update.", req.GetUserId())
		return nil, status.Errorf(codes.InvalidArgument, "No fields to update")
	}

	if err := s.ProfileService.UpdateUserProfileFields(ctx, userId, updates); err != nil {
		log.Printf("ERROR: Failed to update profile for user %s: %v", req.GetUserId(), err)
		return nil, status.Errorf(codes.Internal, "Failed to update profile: %v", err)
	}

	log.Printf("Successfully patched profile for user ID: %s", userId.Hex())
	return &pb.PatchProfileResponse{}, nil
}

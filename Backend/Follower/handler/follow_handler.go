package handler

import (
	"context"

	// Prilagodite putanje vašoj strukturi projekta
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/model"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Follower/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/service"
)

// toProtoUser je helper funkcija koja konvertuje naš interni model.Follow u pb.User.
// U vašem primeru, model se zove Follow, ali predstavlja korisnika.
func toProtoUser(user *model.Follow) *pb.User {
	return &pb.User{
		UserId: user.UserID,
	}
}

// FollowerHandler je gRPC handler za naš Follower servis.
type FollowerHandler struct {
	pb.UnimplementedFollowerServiceServer
	followerService service.FollowService // Koristi interfejs servisnog sloja
}

// NewFollowerHandler kreira novu instancu handlera.
func NewFollowerHandler(followerService service.FollowService) *FollowerHandler {
	return &FollowerHandler{
		followerService: followerService,
	}
}

// FollowUser implementira RPC za praćenje korisnika.
func (h *FollowerHandler) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	// 1. Izvuci podatke iz gRPC zahteva
	followerId := req.GetFollowerId()
	followedId := req.GetFollowedId()

	// 2. Pozovi metodu iz servisnog sloja
	err := h.followerService.FollowUser(ctx, followerId, followedId)
	if err != nil {
		return nil, err
	}

	// 3. Vrati prazan odgovor, jer je status OK dovoljan
	return &pb.FollowUserResponse{}, nil
}
/*
// UnfollowUser implementira RPC za otpraćivanje korisnika.
func (h *FollowerHandler) UnfollowUser(ctx context.Context, req *pb.UnfollowUserRequest) (*pb.UnfollowUserResponse, error) {
	followerId := req.GetFollowerId()
	followedId := req.GetFollowedId()

	// Pretpostavljamo da imate UnfollowUser metodu u servisu
	err := h.followerService.UnfollowUser(ctx, followerId, followedId)
	if err != nil {
		return nil, err
	}

	return &pb.UnfollowUserResponse{}, nil
}

// GetFollowing implementira RPC za dobijanje liste korisnika koje neko prati.
func (h *FollowerHandler) GetFollowing(ctx context.Context, req *pb.GetFollowingRequest) (*pb.GetFollowingResponse, error) {
	// 1. Izvuci ID korisnika iz zahteva
	userId := req.GetUserId()

	// 2. Pozovi servis da dobiješ listu modela
	following, err := h.followerService.GetFollowing(ctx, userId)
	if err != nil {
		return nil, err
	}

	// 3. Konvertuj listu modela u listu proto poruka
	var protoUsers []*pb.User
	for _, user := range following {
		protoUsers = append(protoUsers, toProtoUser(&user))
	}

	// 4. Vrati odgovor
	return &pb.GetFollowingResponse{Users: protoUsers}, nil
}

// GetFollowers implementira RPC za dobijanje liste korisnika koji nekog prate.
func (h *FollowerHandler) GetFollowers(ctx context.Context, req *pb.GetFollowersRequest) (*pb.GetFollowersResponse, error) {
	userId := req.GetUserId()

	followers, err := h.followerService.GetFollowers(ctx, userId)
	if err != nil {
		return nil, err
	}

	var protoUsers []*pb.User
	for _, user := range followers {
		protoUsers = append(protoUsers, toProtoUser(&user))
	}

	return &pb.GetFollowersResponse{Users: protoUsers}, nil
}

// GetFollowRecommendations implementira RPC za dobijanje preporuka.
func (h *FollowerHandler) GetFollowRecommendations(ctx context.Context, req *pb.GetFollowRecommendationsRequest) (*pb.GetFollowRecommendationsResponse, error) {
	userId := req.GetUserId()

	recommendations, err := h.followerService.GetFollowRecommendations(ctx, userId)
	if err != nil {
		return nil, err
	}

	var protoUsers []*pb.User
	for _, user := range recommendations {
		protoUsers = append(protoUsers, toProtoUser(&user))
	}

	return &pb.GetFollowRecommendationsResponse{Users: protoUsers}, nil
}*/
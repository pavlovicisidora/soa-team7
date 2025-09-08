package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// Prilagodite putanje vašoj strukturi projekta
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/model"
	pb "github.com/pavlovicisidora/soa-team7/Backend/Follower/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/service"
)

// toProtoUser je helper funkcija koja konvertuje naš interni model.Follow u pb.User.
// U vašem primeru, model se zove Follow, ali predstavlja korisnika.
func toProtoUser(user *model.User) *pb.User {
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

func (h *FollowerHandler) GetFollowing(ctx context.Context, req *pb.GetFollowingRequest) (*pb.GetFollowingResponse, error) {
	userId := req.GetUserId()

	// 2. Pozovi servis da dobiješ listu modela
	followingUsers, err := h.followerService.GetFollowing(ctx, userId)
	if err != nil {
		return nil, err
	}

	// 3. Konvertuj listu modela u listu proto poruka
	var protoUsers []*pb.User
	for _, userModel := range followingUsers {
		// Koristimo helper funkciju za čistu konverziju
		protoUsers = append(protoUsers, toProtoUser(userModel))
	}

	// 4. Vrati odgovor
	return &pb.GetFollowingResponse{Users: protoUsers}, nil
}


func (h *FollowerHandler) GetFollowRecommendations(ctx context.Context, req *pb.GetFollowRecommendationsRequest) (*pb.GetFollowRecommendationsResponse, error) {
	// 1. Dobijamo ID korisnika iz gRPC zahteva
	userId := req.GetUserId()
	if userId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "User ID cannot be empty")
	}

	// 2. Pozivamo servisnu metodu da dobijemo preporuke
	recommendedUsers, err := h.followerService.GetFollowRecommendations(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get recommendations: %v", err)
	}

	// 3. Konvertujemo rezultat iz modela u proto objekte
	var protoUsers []*pb.User
	for _, userModel := range recommendedUsers {
		// U vašem slučaju helper toProtoUser već postoji i radi ovo, ali ovde je eksplicitno radi jasnoće.
		protoUsers = append(protoUsers, toProtoUser(userModel))
	}

	// 4. Vraćamo odgovor
	return &pb.GetFollowRecommendationsResponse{Users: protoUsers}, nil
}
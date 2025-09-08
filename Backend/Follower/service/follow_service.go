package service

import (
	"context"
	"errors" // Potrebno za definisanje custom grešaka

	"github.com/pavlovicisidora/soa-team7/Backend/Follower/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/repo"
)

// Definišemo custom grešku za lakše rukovanje u handleru.
var ErrCannotFollowSelf = errors.New("user cannot follow themselves")

// UserService definiše interfejs za poslovnu logiku.
type FollowService interface {
	FollowUser(ctx context.Context, followerId string, followedId string) error
	GetFollowing(ctx context.Context, followerId string) ([]*model.User, error)
	GetFollowRecommendations(ctx context.Context, userId string) ([]*model.User, error)
}

// userService je implementacija interfejsa.
type followService struct {
	repo repo.FollowRepository
}

// NewUserService kreira novu instancu servisa.
func NewFollowService(repo repo.FollowRepository) FollowService {
	return &followService{repo: repo}
}

// FollowUser sadrži logiku za praćenje korisnika.
func (s *followService) FollowUser(ctx context.Context, followerId string, followedId string) error {

	if followerId == followedId {
		return ErrCannotFollowSelf
	}

	return s.repo.FollowUser(ctx, followerId, followedId)
}

func (s *followService) GetFollowing(ctx context.Context, followerId string) ([]*model.User, error) {
	// Samo prosleđujemo poziv repozitorijumu.
	// Potpis metode je sada usklađen sa interfejsom i repozitorijumom.
	return s.repo.GetFollowing(ctx, followerId)
}

func (s *followService) GetFollowRecommendations(ctx context.Context, userId string) ([]*model.User, error) {
	return s.repo.GetFollowRecommendations(ctx, userId)
}

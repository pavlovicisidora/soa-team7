package service

import (
	"context"
	"errors" 

	"github.com/pavlovicisidora/soa-team7/Backend/Follower/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/repo"
)


var ErrCannotFollowSelf = errors.New("user cannot follow themselves")


type FollowService interface {
	FollowUser(ctx context.Context, followerId string, followedId string) error
	GetFollowing(ctx context.Context, followerId string) ([]*model.User, error)
	GetFollowRecommendations(ctx context.Context, userId string) ([]*model.User, error)
}


type followService struct {
	repo repo.FollowRepository
}


func NewFollowService(repo repo.FollowRepository) FollowService {
	return &followService{repo: repo}
}


func (s *followService) FollowUser(ctx context.Context, followerId string, followedId string) error {

	if followerId == followedId {
		return ErrCannotFollowSelf
	}

	return s.repo.FollowUser(ctx, followerId, followedId)
}

func (s *followService) GetFollowing(ctx context.Context, followerId string) ([]*model.User, error) {
	return s.repo.GetFollowing(ctx, followerId)
}

func (s *followService) GetFollowRecommendations(ctx context.Context, userId string) ([]*model.User, error) {
	return s.repo.GetFollowRecommendations(ctx, userId)
}

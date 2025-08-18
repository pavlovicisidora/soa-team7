package service

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/model"
	"github.com/pavlovicisidora/soa-team7/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileService struct {
	UserRepo *repo.UserRepository
}

func (s *ProfileService) GetUserProfile(ctx context.Context, userID primitive.ObjectID) (*model.Profile, error) {

	user, err := s.UserRepo.FindUserById(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user.Profile, nil
}

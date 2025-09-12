package service

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/repo"
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

func (s *ProfileService) UpdateUserProfileFields(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}) error {
	return s.UserRepo.UpdateUserProfileFields(ctx, userID, updates)
}

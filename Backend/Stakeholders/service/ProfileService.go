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

//	func (service *ProfileService) FindByUserId(id string) (*model.Profile, error) {
//		profile, err := service.ProfileRepo.FindByUserId(id)
//		if err != nil {
//			return nil, fmt.Errorf("profile item with id not found")
//		}
//		return profile, nil
//	}
func (s *ProfileService) GetUserProfile(ctx context.Context, userID primitive.ObjectID) (*model.Profile, error) {

	user, err := s.UserRepo.FindUserById(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user.Profile, nil
}

package service

import (
	"context"
	"fmt"

	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	UserRepository *repo.UserRepository
}

func (service *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := service.UserRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, err

}
func (service *UserService) Create(ctx context.Context, user *model.User) error {
	existingUser, err := service.UserRepository.FindByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	if existingUser.Username != "" {
		return fmt.Errorf("username already exists")
	}

	existingUser, err = service.UserRepository.FindByMail(ctx, user.Mail)
	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	if existingUser.Mail != "" {
		return fmt.Errorf("there is user with this mail")
	}

	err = service.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) Login(ctx context.Context, username string, password string) (model.User, error) {
	user, err := service.UserRepository.Login(ctx, username, password)

	if err != nil {
		return model.User{}, err
	}
	return user, err

}

func (service *UserService) BlockUser(ctx context.Context, username string) error {
	existingUser, err := service.UserRepository.FindByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	if existingUser.Username != username {
		return fmt.Errorf("this user doesn't exist")
	}

	if existingUser.Role != "VODIC" && existingUser.Role != "TURISTA" {
		return fmt.Errorf("you can't block admins")
	}

	existingUser.Blocked = true

	err = service.UserRepository.UpdateUser(ctx, existingUser)
	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}
	return nil
}
func (service *UserService) FindAllInfo(ctx context.Context, userID string) ([]model.User, error) {
	users, err := service.UserRepository.FindAllInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return users, err

}
func (service *UserService) FindById(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	user, err := service.UserRepository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, err

}

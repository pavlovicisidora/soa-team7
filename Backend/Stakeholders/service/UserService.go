package service

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/model"
	"github.com/pavlovicisidora/soa-team7/repo"
)

type UserService struct {
	UserRepositroy *repo.UserRepository
}

func (service *UserService) Create(ctx context.Context, user *model.User) error {
	err := service.UserRepositroy.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

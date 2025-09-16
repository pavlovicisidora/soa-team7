package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/repo"
	"github.com/pavlovicisidora/soa-team7/Backend/saga"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserService struct {
	UserRepository *repo.UserRepository
	NatsConn       *nats.Conn
}

func (service *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := service.UserRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, err

}
func (service *UserService) Create(ctx context.Context, user *model.User) error {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "user.service.create")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.username", user.Username),
		attribute.String("user.role", user.Role),
	)
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
		span.RecordError(err)
		return err
	}
	return nil
}

func (service *UserService) Login(ctx context.Context, username string, password string) (model.User, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "user.service.login")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.username", username),
	)

	user, err := service.UserRepository.Login(ctx, username, password)

	if err != nil {
		span.RecordError(err)
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

	event := saga.UserBlockedEvent{
		UserID: existingUser.ID.Hex(),
	}
	eventData, _ := json.Marshal(event)

	log.Printf("Publishing event to subject: %s for UserID: %s", saga.UserBlockedSubject, event.UserID)
	if err := service.NatsConn.Publish(saga.UserBlockedSubject, eventData); err != nil {
		log.Printf("Failed to publish NATS event, compensating... Error: %v", err)
		existingUser.Blocked = false
		_ = service.UserRepository.UpdateUser(ctx, existingUser)
		return fmt.Errorf("failed to publish saga event")
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

func (service *UserService) UpdateUserPosition(ctx context.Context, userID string, lat, long float64) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid userID format: %v", err)
	}
	return service.UserRepository.UpdateUserPosition(ctx, userObjectID, lat, long)
}

func (service *UserService) HandleBlockUserCompensation(ctx context.Context, userID string) error {
	log.Printf("Executing compensation for UserID: %s", userID)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid userID format for compensation: %v", err)
	}

	user, err := service.UserRepository.FindUserById(ctx, userObjectID)
	if err != nil {
		return fmt.Errorf("could not find user for compensation: %v", err)
	}

	// Vraćanje stanja
	user.Blocked = false
	if err := service.UserRepository.UpdateUser(ctx, *user); err != nil {
		return fmt.Errorf("failed to update user during compensation: %v", err)
	}

	log.Printf("Compensation successful for UserID: %s", userID)
	return nil
}

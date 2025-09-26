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
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.GetAllUsers")
	defer span.End()
	users, err := service.UserRepository.GetAllUsers(ctx)
	if err != nil {
		span.RecordError(err)
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
	log.Printf("SERVICE: Checking if username '%s' exists.", user.Username)
	existingUser, err := service.UserRepository.FindByUsername(ctx, user.Username)
	if err != nil {
		log.Printf("ERROR: DB error while checking username '%s': %v", user.Username, err)
		return fmt.Errorf("DB error: %v", err)
	}

	if existingUser.Username != "" {
		log.Printf("WARN: Username '%s' already exists.", user.Username)
		return fmt.Errorf("username already exists")
	}

	log.Printf("SERVICE: Checking if mail '%s' exists.", user.Mail)
	existingUser, err = service.UserRepository.FindByMail(ctx, user.Mail)
	if err != nil {
		log.Printf("ERROR: DB error while checking mail '%s': %v", user.Mail, err)
		return fmt.Errorf("DB error: %v", err)
	}

	if existingUser.Mail != "" {
		log.Printf("WARN: Mail '%s' already exists for another user.", user.Mail)
		return fmt.Errorf("there is user with this mail")
	}

	log.Printf("SERVICE: Creating user '%s' in repository.", user.Username)
	err = service.UserRepository.CreateUser(ctx, user)
	if err != nil {
		span.RecordError(err)
		log.Printf("ERROR: Failed to create user '%s' in repository: %v", user.Username, err)
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
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.BlockUser")
	defer span.End()
	span.SetAttributes(attribute.String("user.username", username))

	log.Printf("SERVICE: Attempting to block user '%s'.", username)
	existingUser, err := service.UserRepository.FindByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		log.Printf("ERROR: DB error while finding user '%s' to block: %v", username, err)
		return fmt.Errorf("DB error: %v", err)
	}

	if existingUser.Username != username {
		log.Printf("WARN: User '%s' not found for blocking.", username)
		return fmt.Errorf("this user doesn't exist")
	}

	if existingUser.Role != "VODIC" && existingUser.Role != "TURISTA" {
		log.Printf("WARN: Attempted to block an admin user: %s", username)
		return fmt.Errorf("you can't block admins")
	}

	existingUser.Blocked = true

	log.Printf("SERVICE: Updating user '%s' to blocked=true in repository.", username)
	err = service.UserRepository.UpdateUser(ctx, existingUser)
	if err != nil {
		log.Printf("ERROR: DB error while updating user '%s' to blocked: %v", username, err)
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

	log.Printf("SERVICE: Successfully blocked user '%s' and published SAGA event.", username)
	return nil
}
func (service *UserService) FindAllInfo(ctx context.Context, userID string) ([]model.User, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.FindAllInfo")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", userID))
	users, err := service.UserRepository.FindAllInfo(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return users, err

}
func (service *UserService) FindById(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.FindById")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", id.Hex()))
	user, err := service.UserRepository.FindUserById(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return user, err

}

func (service *UserService) UpdateUserPosition(ctx context.Context, userID string, lat, long float64) error {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.UpdateUserPosition")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", userID))
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("invalid userID format: %v", err)
	}
	return service.UserRepository.UpdateUserPosition(ctx, userObjectID, lat, long)
}

func (service *UserService) HandleBlockUserCompensation(userID string) error {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(context.Background(), "service.HandleBlockUserCompensation") // Kreira novi, "root" span
	defer span.End()
	span.SetAttributes(attribute.String("user.id", userID))
	log.Printf("Executing compensation for UserID: %s", userID)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("SAGA_COMPENSATION_ERROR: Invalid userID format: %v", err)
		return fmt.Errorf("invalid userID format for compensation: %v", err)
	}

	user, err := service.UserRepository.FindUserById(ctx, userObjectID)
	if err != nil {
		log.Printf("SAGA_COMPENSATION_ERROR: Could not find user for compensation: %v", err)
		return fmt.Errorf("could not find user for compensation: %v", err)
	}

	// Vraćanje stanja
	user.Blocked = false
	log.Printf("SAGA_COMPENSATION: Setting user %s (ID: %s) back to blocked=false.", user.Username, userID)
	if err := service.UserRepository.UpdateUser(ctx, *user); err != nil {
		log.Printf("SAGA_COMPENSATION_ERROR: Failed to update user during compensation: %v", err)
		return fmt.Errorf("failed to update user during compensation: %v", err)
	}

	log.Printf("Compensation successful for UserID: %s", userID)
	return nil
}

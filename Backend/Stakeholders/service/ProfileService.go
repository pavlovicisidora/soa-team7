package service

import (
	"context"
	"log"

	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ProfileService struct {
	UserRepo *repo.UserRepository
}

func (s *ProfileService) GetUserProfile(ctx context.Context, userID primitive.ObjectID) (*model.Profile, error) {

	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.GetUserProfile")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", userID.Hex()))
	log.Printf("SERVICE: Getting user profile for ID: %s", userID.Hex())
	user, err := s.UserRepo.FindUserById(ctx, userID)
	if err != nil {
		span.RecordError(err)
		log.Printf("ERROR: Failed to find user by ID %s to get profile: %v", userID.Hex(), err)
		return nil, err
	}
	return &user.Profile, nil
}

func (s *ProfileService) UpdateUserProfileFields(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}) error {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.UpdateUserProfileFields")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", userID.Hex()))
	log.Printf("SERVICE: Updating profile for user ID: %s with %d fields.", userID.Hex(), len(updates))
	err := s.UserRepo.UpdateUserProfileFields(ctx, userID, updates)
	if err != nil {
		span.RecordError(err)
		log.Printf("ERROR: Failed to update profile fields for user ID %s: %v", userID.Hex(), err)
	}
	return err
}

package service

import (
	"context"

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
	user, err := s.UserRepo.FindUserById(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return &user.Profile, nil
}

func (s *ProfileService) UpdateUserProfileFields(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}) error {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.UpdateUserProfileFields")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", userID.Hex()))
	err := s.UserRepo.UpdateUserProfileFields(ctx, userID, updates)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

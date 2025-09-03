package service

import (
	"context"
	"time"

	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	AddComment(ctx context.Context, blogID, authorID, text string) (*model.Comment, error)
	GetComments(ctx context.Context, blogID string) ([]model.Comment, error)
	UpdateComment(ctx context.Context, commentID, newText string) (*model.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
}

type commentService struct {
	repo repo.CommentRepository
}

func NewCommentService(r repo.CommentRepository) CommentService {
	return &commentService{repo: r}
}

// AddComment kreira novi komentar
func (s *commentService) AddComment(ctx context.Context, blogID, authorID, text string) (*model.Comment, error) {
	blogOID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, err
	}

	comment := &model.Comment{
		ID:             primitive.NewObjectID(),
		BlogID:         blogOID,
		AuthorID:       authorID,
		Text:           text,
		CreatedAt:      time.Now(),
	}

	_, err = s.repo.AddComment(ctx, blogOID, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetComments vraća sve komentare za dati blog
func (s *commentService) GetComments(ctx context.Context, blogID string) ([]model.Comment, error) {
	blogOID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetComments(ctx, blogOID)
}

// UpdateComment menja tekst komentara po njegovom ID
func (s *commentService) UpdateComment(ctx context.Context, commentID, newText string) (*model.Comment, error) {
	commentOID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdateComment(ctx, commentOID, newText)
}

// DeleteComment briše komentar po njegovom ID
func (s *commentService) DeleteComment(ctx context.Context, commentID string) error {
	commentOID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	return s.repo.DeleteComment(ctx, commentOID)
}

package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type BlogService interface {
	CreateBlog(ctx context.Context, title, content string, images []model.Image, userID string) (*model.Blog, error)
	GetAllBlogs(ctx context.Context, authorIDs []string) ([]model.Blog, error)
	GetBlogByID(ctx context.Context, id string) (*model.Blog, error)
	LikeBlog(ctx context.Context, blogID, userID string) (*model.Blog, error)
	UnlikeBlog(ctx context.Context, blogID, userID string) (*model.Blog, error)
	HandleUserBlocked(ctx context.Context, userID string) error
}

type blogService struct {
	blogRepo repo.BlogRepository
	NatsConn *nats.Conn
}

func NewBlogService(blogRepo repo.BlogRepository, nc *nats.Conn) BlogService {
	return &blogService{
		blogRepo: blogRepo,
		NatsConn: nc,
	}
}

func (s *blogService) CreateBlog(ctx context.Context, title, content string, images []model.Image, userID string) (*model.Blog, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.CreateBlog")
	defer span.End()

	log.Printf("SERVICE: Creating new blog with title '%s' for user %s", title, userID)

	newBlog := &model.Blog{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		Images:    images,
		UserID:    userID,
	}

	insertedID, err := s.blogRepo.CreateBlog(ctx, newBlog)
	if err != nil {
		return nil, handleServiceError(span, "Failed to create blog in repository", err)
	}

	newBlog.ID = *insertedID
	log.Printf("SERVICE: Blog created successfully with ID: %s", newBlog.ID.Hex())
	return newBlog, nil
}

func (s *blogService) GetAllBlogs(ctx context.Context, authorIDs []string) ([]model.Blog, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.GetAllBlogs")
	defer span.End()

	log.Printf("SERVICE: Getting all blogs for %d authors", len(authorIDs))
	blogs, err := s.blogRepo.GetBlogs(ctx, authorIDs)
	if err != nil {
		return nil, handleServiceError(span, "Failed to get all blogs from repository", err)
	}
	return blogs, nil

}

func (s *blogService) GetBlogByID(ctx context.Context, id string) (*model.Blog, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.GetBlogByID")
	defer span.End()

	log.Printf("SERVICE: Getting blog by ID: %s", id)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, handleServiceError(span, fmt.Sprintf("Invalid blog ID format: %s", id), err)
	}

	blog, err := s.blogRepo.GetBlogByID(ctx, objID)
	if err != nil {
		return nil, handleServiceError(span, fmt.Sprintf("Failed to get blog by ID %s from repository", id), err)
	}
	return blog, nil
}

func (s *blogService) LikeBlog(ctx context.Context, blogID, userID string) (*model.Blog, error) {
	objBlogID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, err
	}
	if err := s.blogRepo.LikeBlog(ctx, objBlogID, userID); err != nil {
		return nil, err
	}
	return s.blogRepo.GetBlogByID(ctx, objBlogID)
}

func (s *blogService) UnlikeBlog(ctx context.Context, blogID, userID string) (*model.Blog, error) {
	objBlogID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, err
	}
	if err := s.blogRepo.UnlikeBlog(ctx, objBlogID, userID); err != nil {
		return nil, err
	}
	return s.blogRepo.GetBlogByID(ctx, objBlogID)
}
func (s *blogService) HandleUserBlocked(ctx context.Context, userID string) error {
	// return fmt.Errorf("simulirana greška za testiranje SAGA rollback-a")
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "service.HandleUserBlocked_SAGA")
	defer span.End()

	log.Printf("SAGA_STEP: Handling UserBlocked event for user ID: %s", userID)
	err := s.blogRepo.UpdateBlogsOnUserStatusChange(ctx, userID, true)
	if err != nil {
		return handleServiceError(span, fmt.Sprintf("Failed to update blogs for blocked user %s", userID), err)
	}
	log.Printf("SAGA_STEP: Successfully updated blogs for blocked user ID: %s", userID)
	return err
}
func handleServiceError(span trace.Span, message string, err error) error {
	log.Printf("ERROR: "+message+": %v", err)
	span.RecordError(err)
	span.SetStatus(codes.Error, message)
	return fmt.Errorf(message+": %w", err)
}

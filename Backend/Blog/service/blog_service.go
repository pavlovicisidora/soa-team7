package service

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	newBlog := &model.Blog{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		Images:    images,
		UserID:    userID,
	}

	insertedID, err := s.blogRepo.CreateBlog(ctx, newBlog)
	if err != nil {
		return nil, err
	}

	newBlog.ID = *insertedID
	return newBlog, nil
}

func (s *blogService) GetAllBlogs(ctx context.Context, authorIDs []string) ([]model.Blog, error) {
	return s.blogRepo.GetBlogs(ctx, authorIDs)
}

func (s *blogService) GetBlogByID(ctx context.Context, id string) (*model.Blog, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.blogRepo.GetBlogByID(ctx, objID)
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
	err := s.blogRepo.UpdateBlogsOnUserStatusChange(ctx, userID, true)
	return err
}

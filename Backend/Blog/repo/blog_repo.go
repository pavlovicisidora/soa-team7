package repo

import (
	"context"
	"fmt"
	"log"

	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *model.Blog) (*primitive.ObjectID, error)
	GetBlogs(ctx context.Context, authorIDs []string) ([]model.Blog, error)
	GetBlogByID(ctx context.Context, id primitive.ObjectID) (*model.Blog, error)
	LikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error
	UnlikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error
	UpdateBlogsOnUserStatusChange(ctx context.Context, userID string, isBlocked bool) error
}

type blogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(collection *mongo.Collection) BlogRepository {
	return &blogRepository{
		collection: collection,
	}
}

func (r *blogRepository) CreateBlog(ctx context.Context, blog *model.Blog) (*primitive.ObjectID, error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "repo.CreateBlog")
	defer span.End()
	span.SetAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.statement", "InsertOne"),
	)

	log.Printf("REPO: Inserting new blog into database.")
	result, err := r.collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, handleRepoError(span, "DB InsertOne failed", err)
	}
	id := result.InsertedID.(primitive.ObjectID)
	log.Printf("REPO: Blog inserted with ID: %s", id.Hex())
	return &id, nil
}

func (r *blogRepository) GetBlogs(ctx context.Context, authorIDs []string) ([]model.Blog, error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "repo.GetBlogs")
	defer span.End()
	span.SetAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.statement", "Find"),
	)
	if len(authorIDs) == 0 {
		log.Printf("REPO: GetBlogs called with no author IDs, returning empty list.")
		return []model.Blog{}, nil
	}

	filter := bson.M{
		"user_id":        bson.M{"$in": authorIDs},
		"author_blocked": bson.M{"$ne": true},
	}

	var blogs []model.Blog
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, handleRepoError(span, "DB Find failed", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, handleRepoError(span, "DB cursor decoding failed", err)
	}
	log.Printf("REPO: Fetched %d blogs from database.", len(blogs))
	return blogs, nil
}

func (r *blogRepository) GetBlogByID(ctx context.Context, id primitive.ObjectID) (*model.Blog, error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "repo.GetBlogByID")
	defer span.End()
	span.SetAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.statement", "FindOne"),
		attribute.String("blog.id", id.Hex()),
	)
	log.Printf("REPO: Getting blog by ID: %s", id.Hex())
	var blog model.Blog
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		return nil, handleRepoError(span, fmt.Sprintf("DB FindOne failed for blog ID %s", id.Hex()), err)
	}
	log.Printf("REPO: Successfully fetched blog by ID: %s", id.Hex())
	return &blog, nil
}

func (r *blogRepository) LikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	filter := bson.M{"_id": blogID}
	update := bson.M{"$addToSet": bson.M{"liked_by": userID}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *blogRepository) UnlikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	filter := bson.M{"_id": blogID}
	update := bson.M{"$pull": bson.M{"liked_by": userID}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
func (r *blogRepository) UpdateBlogsOnUserStatusChange(ctx context.Context, userID string, isBlocked bool) error {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "repo.UpdateBlogsOnUserStatusChange")
	defer span.End()
	span.SetAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.statement", "UpdateMany"),
		attribute.String("user.id", userID),
		attribute.Bool("user.isBlocked", isBlocked),
	)
	log.Printf("REPO: Updating blogs for user %s, setting author_blocked to %t", userID, isBlocked)
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"author_blocked": isBlocked}}

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return handleRepoError(span, fmt.Sprintf("DB UpdateMany failed for user status change %s", userID), err)
	}
	log.Printf("REPO: Successfully updated %d blogs for user %s", result.ModifiedCount, userID)
	return nil
}

func handleRepoError(span trace.Span, message string, err error) error {
	log.Printf("ERROR: "+message+": %v", err)
	span.RecordError(err)
	span.SetStatus(codes.Error, message)
	return fmt.Errorf(message+": %w", err)
}

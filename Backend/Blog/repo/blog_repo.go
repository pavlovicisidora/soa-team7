package repo

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *model.Blog) (*primitive.ObjectID, error)
	GetBlogs(ctx context.Context, authorIDs []string) ([]model.Blog, error)
	GetBlogByID(ctx context.Context, id primitive.ObjectID) (*model.Blog, error)
	LikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error
	UnlikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error
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
	result, err := r.collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (r *blogRepository) GetBlogs(ctx context.Context, authorIDs []string) ([]model.Blog, error) {
	// Ako je lista ID-jeva prazna, nema potrebe za upitom, vraćamo prazan slice
	if len(authorIDs) == 0 {
		return []model.Blog{}, nil
	}

	// Koristimo $in operator da nađemo sve blogove gde se user_id poklapa sa bilo kojim ID-jem iz liste
	filter := bson.M{"user_id": bson.M{"$in": authorIDs}}

	var blogs []model.Blog
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *blogRepository) GetBlogByID(ctx context.Context, id primitive.ObjectID) (*model.Blog, error) {
	var blog model.Blog
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		return nil, err
	}
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

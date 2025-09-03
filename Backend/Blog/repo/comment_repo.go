package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository interface {
	AddComment(ctx context.Context, blogID primitive.ObjectID, comment *model.Comment) (*primitive.ObjectID, error)
	GetComments(ctx context.Context, blogID primitive.ObjectID) ([]model.Comment, error)
	UpdateComment(ctx context.Context, commentID primitive.ObjectID, newText string) (*model.Comment, error)
	DeleteComment(ctx context.Context, commentID primitive.ObjectID) error
}

type commentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(collection *mongo.Collection) CommentRepository {
	return &commentRepository{collection: collection}
}

func (r *commentRepository) AddComment(ctx context.Context, blogID primitive.ObjectID, comment *model.Comment) (*primitive.ObjectID, error) {
	if comment.ID.IsZero() {
		comment.ID = primitive.NewObjectID()
	}
	comment.BlogID = blogID
	comment.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}
	return &comment.ID, nil
}

func (r *commentRepository) GetComments(ctx context.Context, blogID primitive.ObjectID) ([]model.Comment, error) {
	filter := bson.M{"blog_id": blogID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var comments []model.Comment
	for cur.Next(ctx) {
		var c model.Comment
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *commentRepository) UpdateComment(ctx context.Context, commentID primitive.ObjectID, newText string) (*model.Comment, error) {
	filter := bson.M{"_id": commentID}
	update := bson.M{
		"$set": bson.M{
			"text":       newText,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("comment not found")
	}

	var updated model.Comment
	if err := r.collection.FindOne(ctx, filter).Decode(&updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, commentID primitive.ObjectID) error {
	filter := bson.M{"_id": commentID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("comment not found")
	}
	return nil
}

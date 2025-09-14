package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	URL string `bson:"url" json:"url"`
}

type Blog struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title         string             `bson:"title" json:"title"`
	Content       string             `bson:"content" json:"content"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	Images        []Image            `bson:"images" json:"images,omitempty"`
	UserID        string             `bson:"user_id" json:"user_id"`
	LikedBy       []string           `bson:"liked_by,omitempty" json:"liked_by"`
	AuthorBlocked bool               `bson:"author_blocked,omitempty" json:"author_blocked,omitempty"`
}

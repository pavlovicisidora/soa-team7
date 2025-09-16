package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItem struct {
	TourID   int     `bson:"tour_id" json:"tour_id"`
	TourName string  `bson:"tour_name" json:"tour_name"`
	Price    float64 `bson:"price" json:"price"`
}

type ShoppingCart struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     string             `bson:"user_id" json:"user_id"`
	Items      []OrderItem        `bson:"items" json:"items"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
}

type TourPurchaseToken struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID string             `bson:"user_id" json:"user_id"`
	TourID int                `bson:"tour_id" json:"tour_id"`
}

package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Mail     string             `json:"mail" bson:"mail"`
	Role     string             `json:"role" bson:"role"`
	Blocked  bool               `json:"blocked" bson:"blocked"`
	Profile  Profile            `json:"profile" bson:"profile"`
}

package repo

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *UserRepository {
	return &UserRepository{
		client:         client,
		dbName:         dbName,
		collectionName: collectionName,
	}
}


func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

package repo

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/model"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	// Empty filter = get all documents
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)    //Ovo pomaze da na dodje do leaka, cursor predstavlja otvoreni reyultat upita pa ga zatvarmao da 

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

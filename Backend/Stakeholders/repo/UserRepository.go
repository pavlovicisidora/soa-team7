package repo

import (
	"context"
	"fmt"

	"github.com/pavlovicisidora/soa-team7/Backend/Stakeholders/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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
	defer cursor.Close(ctx) //Ovo pomaze da na dodje do leaka, cursor predstavlja otvoreni reyultat upita pa ga zatvarmao da

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user model.User) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	filter := bson.M{"username": user.Username}

	update := bson.M{
		"$set": bson.M{
			"mail":     user.Mail,
			"role":     user.Role,
			"blocked":  user.Blocked,
			"password": user.Password,
			"profile":  user.Profile,
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// ako nije pronađen nijedan dokument
	if result.MatchedCount == 0 {
		return fmt.Errorf("user with username %s not found", user.Username)
	}

	return nil

}

func (r *UserRepository) Login(ctx context.Context, username string, password string) (model.User, error) {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.User{}, fmt.Errorf("this username doesn't exist")
		}
		return model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.User{}, fmt.Errorf("wrong password")
	}

	return user, nil

}
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (model.User, error) {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.User{}, nil
		}
		return model.User{}, err
	}
	return user, err

}

func (r *UserRepository) FindByMail(ctx context.Context, mail string) (model.User, error) {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	var user model.User
	err := collection.FindOne(ctx, bson.M{"mail": mail}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.User{}, nil
		}
		return model.User{}, err
	}
	return user, err
}
func (r *UserRepository) FindAllInfo(ctx context.Context, userID string) ([]model.User, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID format: %v", err)
	}
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	projection := bson.M{
		"username": 1,
		"mail":     1,
		"role":     1,
		"_id":      0,
		"blocked":  1,
	}
	findOptions := options.Find().SetProjection(projection)
	filter := bson.M{"_id": bson.M{"$ne": userObjectID}}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
func (r *UserRepository) FindUserById(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	var user model.User

	filter := bson.M{"_id": id}

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

/*
func (r *UserRepository) UpdateUserProfileById(ctx context.Context, id primitive.ObjectID, profile model.Profile) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"profile": profile,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user with id %s not found", id.Hex())
	}

	return nil
}*/

func (r *UserRepository) UpdateUserProfileFields(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	update := bson.M{"$set": bson.M{}}
	for key, value := range updates {
		update["$set"].(bson.M)["profile."+key] = value
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user with id %s not found", id.Hex())
	}

	return nil
}

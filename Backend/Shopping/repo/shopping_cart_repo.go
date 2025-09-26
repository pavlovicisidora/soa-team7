package repo

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/Backend/Shopping/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShoppingCartRepository struct {
	cartsCollection  *mongo.Collection
	tokensCollection *mongo.Collection
}

func NewShoppingCartRepository(db *mongo.Database) *ShoppingCartRepository {
	return &ShoppingCartRepository{
		cartsCollection:  db.Collection("shopping_carts"),
		tokensCollection: db.Collection("purchase_tokens"),
	}
}

func (r *ShoppingCartRepository) GetCartByUserID(ctx context.Context, userID string) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart
	filter := bson.M{"user_id": userID}
	err := r.cartsCollection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *ShoppingCartRepository) UpsertCart(ctx context.Context, cart *model.ShoppingCart) error {
	filter := bson.M{"user_id": cart.UserID}
	update := bson.M{"$set": cart}
	opts := options.Update().SetUpsert(true)

	_, err := r.cartsCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *ShoppingCartRepository) CreatePurchaseTokens(ctx context.Context, tokens []interface{}) error {
	if len(tokens) == 0 {
		return nil
	}
	_, err := r.tokensCollection.InsertMany(ctx, tokens)
	return err
}

func (r *ShoppingCartRepository) DeleteCart(ctx context.Context, userID string) error {
	filter := bson.M{"user_id": userID}
	_, err := r.cartsCollection.DeleteOne(ctx, filter)
	return err
}

func (r *ShoppingCartRepository) HasToken(ctx context.Context, userID string, tourID int) (bool, error) {
	filter := bson.M{"user_id": userID, "tour_id": tourID}
	count, err := r.tokensCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

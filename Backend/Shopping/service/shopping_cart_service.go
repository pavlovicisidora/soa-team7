package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	tour_proto "github.com/pavlovicisidora/soa-team7/Backend/APIGateway/proto"
	"github.com/pavlovicisidora/soa-team7/Backend/Shopping/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Shopping/repo"
)

type ShoppingCartService struct {
	repo       *repo.ShoppingCartRepository
	tourClient tour_proto.TourGrpcServiceClient
}

func NewShoppingCartService(repo *repo.ShoppingCartRepository, tourClient tour_proto.TourGrpcServiceClient) *ShoppingCartService {
	return &ShoppingCartService{repo: repo, tourClient: tourClient}
}

func (s *ShoppingCartService) recalculateTotalPrice(cart *model.ShoppingCart) {
	cart.TotalPrice = 0
	for _, item := range cart.Items {
		cart.TotalPrice += item.Price
	}
}

func (s *ShoppingCartService) GetCart(ctx context.Context, userID string) (*model.ShoppingCart, error) {
	cart, err := s.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return &model.ShoppingCart{UserID: userID, Items: []model.OrderItem{}, TotalPrice: 0}, nil
	}
	return cart, nil
}

func (s *ShoppingCartService) AddItemToCart(ctx context.Context, userID string, tourID int) (*model.ShoppingCart, error) {
	tourResponse, err := s.tourClient.GetTourById(ctx, &tour_proto.GetTourByIdRequest{TourId: int32(tourID)})
	if err != nil {
		return nil, fmt.Errorf("could not get tour details: %w", err)
	}
	tour := tourResponse.GetTour()
	if tour.Status != "PUBLISHED" {
		return nil, errors.New("only published tours can be added to the cart")
	}

	cart, err := s.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		cart = &model.ShoppingCart{UserID: userID, Items: []model.OrderItem{}}
	}

	for _, item := range cart.Items {
		if item.TourID == tourID {
			return cart, nil
		}
	}

	newItem := model.OrderItem{
		TourID:   tourID,
		TourName: tour.Name,
		Price:    tour.Price,
	}
	cart.Items = append(cart.Items, newItem)

	s.recalculateTotalPrice(cart)

	if err := s.repo.UpsertCart(ctx, cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *ShoppingCartService) RemoveItemFromCart(ctx context.Context, userID string, tourID int) (*model.ShoppingCart, error) {
	cart, err := s.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return &model.ShoppingCart{UserID: userID}, nil
	}

	itemFound := false
	var updatedItems []model.OrderItem
	for _, item := range cart.Items {
		if item.TourID == tourID {
			itemFound = true
		} else {
			updatedItems = append(updatedItems, item)
		}
	}

	if !itemFound {
		return cart, nil
	}

	cart.Items = updatedItems
	s.recalculateTotalPrice(cart)

	if err := s.repo.UpsertCart(ctx, cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *ShoppingCartService) Checkout(ctx context.Context, userID string) ([]model.TourPurchaseToken, error) {
	log.Printf("SERVICE: Starting checkout for userID=[%s]", userID)

	cart, err := s.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get cart for user %s: %v", userID, err)
		return nil, err
	}
	if cart == nil || len(cart.Items) == 0 {
		log.Println("ERROR: User tried to checkout with an empty cart.")
		return nil, errors.New("shopping cart is empty")
	}

	log.Printf("SERVICE: Found cart with %d items. Total price: %f", len(cart.Items), cart.TotalPrice)

	var tokensToCreate []interface{}
	var createdTokens []model.TourPurchaseToken
	for _, item := range cart.Items {
		token := model.TourPurchaseToken{
			UserID: userID,
			TourID: item.TourID,
		}
		tokensToCreate = append(tokensToCreate, token)
		createdTokens = append(createdTokens, token)
	}

	log.Printf("SERVICE: Prepared %d tokens to be created.", len(tokensToCreate))

	if err := s.repo.CreatePurchaseTokens(ctx, tokensToCreate); err != nil {
		log.Printf("ERROR: Failed to create purchase tokens in database: %v", err)
		return nil, fmt.Errorf("database error during token creation: %w", err)
	}

	log.Println("SERVICE: Successfully created purchase tokens.")

	if err := s.repo.DeleteCart(ctx, userID); err != nil {
		log.Printf("WARNING: Failed to delete cart for user %s after checkout: %v", userID, err)
	} else {
		log.Println("SERVICE: Successfully deleted cart after checkout.")
	}

	log.Println("SERVICE: Checkout completed successfully.")
	return createdTokens, nil
}

func (s *ShoppingCartService) CheckToken(ctx context.Context, userID string, tourID int) (bool, error) {
	return s.repo.HasToken(ctx, userID, tourID)
}

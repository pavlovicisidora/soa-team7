package repo

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/pavlovicisidora/soa-team7/Backend/Follower/model"
)

// UserRepository definiše interfejs za operacije sa korisnicima u bazi.
type FollowRepository interface {
	FollowUser(ctx context.Context, followerId string, followedId string) error
	GetFollowing(ctx context.Context, followerId string) ([]*model.User, error)
	GetFollowRecommendations(ctx context.Context, userId string) ([]*model.User, error)
}

// neo4jUserRepository je implementacija interfejsa za Neo4j bazu.
type neo4jFollowRepository struct {
	driver neo4j.DriverWithContext
}

// NewNeo4jUserRepository kreira novu instancu repository-ja.
func NewNeo4jFollowRepository(driver neo4j.DriverWithContext) FollowRepository {
	return &neo4jFollowRepository{driver: driver}
}

func (r *neo4jFollowRepository) FollowUser(ctx context.Context, followerId string, followedId string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Koristimo ExecuteWrite za transakcione operacije pisanja.
	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
               
                MERGE (follower:User {userId: $followerId})

                
                MERGE (followed:User {userId: $followedId})

                
                MERGE (follower)-[:FOLLOWS]->(followed)
            `
			parameters := map[string]interface{}{
				"followerId": followerId,
				"followedId": followedId,
			}

			_, err := tx.Run(ctx, query, parameters)

			return nil, err
		})

	if err != nil {
		return fmt.Errorf("could not execute follow user query: %w", err)
	}

	return nil
}

func (r *neo4jFollowRepository) GetFollowing(ctx context.Context, followerId string) ([]*model.User, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
                MATCH (follower:User {userId: $followerId})-[:FOLLOWS]->(followed:User)
                RETURN followed.userId AS followedUserId
            `
			parameters := map[string]interface{}{
				"followerId": followerId,
			}

			res, err := tx.Run(ctx, query, parameters)
			if err != nil {
				return nil, err
			}

			// <-- LOGIKA IZRADE MODELA
			var followingUsers []*model.User // Kreiramo slice pointera na model
			for res.Next(ctx) {
				record := res.Record()
				id, ok := record.Get("followedUserId")
				if ok {
					// Za svaki rezultat, kreiramo instancu našeg modela
					user := &model.User{
						UserID: id.(string),
					}
					followingUsers = append(followingUsers, user)
				}
			}

			if err = res.Err(); err != nil {
				return nil, err
			}

			return followingUsers, nil
		})

	if err != nil {
		return nil, fmt.Errorf("could not get following users: %w", err)
	}

	if result == nil {
		return []*model.User{}, nil // Vraćamo prazan slice umesto nil
	}

	// Kastujemo rezultat u odgovarajući tip
	return result.([]*model.User), nil
}

func (r *neo4jFollowRepository) GetFollowRecommendations(ctx context.Context, userId string) ([]*model.User, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
                MATCH (me:User {userId: $userId})-[:FOLLOWS]->(:User)-[:FOLLOWS]->(recommendation:User)
                WHERE NOT (me)-[:FOLLOWS]->(recommendation) AND me <> recommendation
                RETURN DISTINCT recommendation.userId AS recommendedUserId
                LIMIT 10
            `
			parameters := map[string]interface{}{
				"userId": userId,
			}

			res, err := tx.Run(ctx, query, parameters)
			if err != nil {
				return nil, err
			}

			var recommendedUsers []*model.User
			for res.Next(ctx) {
				record := res.Record()
				id, ok := record.Get("recommendedUserId")
				if ok {
					user := &model.User{
						UserID: id.(string),
					}
					recommendedUsers = append(recommendedUsers, user)
				}
			}

			if err = res.Err(); err != nil {
				return nil, err
			}

			return recommendedUsers, nil
		})

	if err != nil {
		return nil, fmt.Errorf("could not get follow recommendations: %w", err)
	}

	if result == nil {
		return []*model.User{}, nil
	}

	return result.([]*model.User), nil
}

package repo

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	//"github.com/pavlovicisidora/soa-team7/Backend/Follower/model"
)

// UserRepository definiše interfejs za operacije sa korisnicima u bazi.
type FollowRepository interface {
	FollowUser(ctx context.Context, followerId string, followedId string) error
	// Ovde će doći i druge metode kao što su GetFollowRecommendations, UnfollowUser, itd.
}

// neo4jUserRepository je implementacija interfejsa za Neo4j bazu.
type neo4jFollowRepository struct {
	driver neo4j.DriverWithContext
}

// NewNeo4jUserRepository kreira novu instancu repository-ja.
func NewNeo4jFollowRepository(driver neo4j.DriverWithContext) FollowRepository {
	return &neo4jFollowRepository{driver: driver}
}

// FollowUser kreira :FOLLOWS vezu između dva korisnika.
// Ako korisnici ne postoje, biće kreirani (samo sa userId).
func (r *neo4jFollowRepository) FollowUser(ctx context.Context, followerId string, followedId string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Koristimo ExecuteWrite za transakcione operacije pisanja.
	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
                // Pronađi ili kreiraj čvor za korisnika koji prati (follower)
                MERGE (follower:User {userId: $followerId})

                // Pronađi ili kreiraj čvor za korisnika koji je praćen (followed)
                MERGE (followed:User {userId: $followedId})

                // Spoji ih vezom :FOLLOWS, ako već ne postoji
                MERGE (follower)-[:FOLLOWS]->(followed)
            `
			parameters := map[string]interface{}{
				"followerId": followerId,
				"followedId": followedId,
			}

			// Izvrši upit
			_, err := tx.Run(ctx, query, parameters)

			// Vraćamo grešku ako je do nje došlo unutar transakcije
			return nil, err
		})

	// Vraćamo grešku ako je došlo do problema sa sesijom ili konekcijom
	if err != nil {
		return fmt.Errorf("could not execute follow user query: %w", err)
	}

	return nil
}

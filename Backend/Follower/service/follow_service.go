package service

import (
	"context"
	"errors" // Potrebno za definisanje custom grešaka

	"github.com/pavlovicisidora/soa-team7/Backend/Follower/repo"
)

// Definišemo custom grešku za lakše rukovanje u handleru.
var ErrCannotFollowSelf = errors.New("user cannot follow themselves")

// UserService definiše interfejs za poslovnu logiku.
type FollowService interface {
	FollowUser(ctx context.Context, followerId string, followedId string) error
}

// userService je implementacija interfejsa.
type followService struct {
	repo repo.FollowRepository
}

// NewUserService kreira novu instancu servisa.
func NewFollowService(repo repo.FollowRepository) FollowService {
	return &followService{repo: repo}
}

// FollowUser sadrži logiku za praćenje korisnika.
func (s *followService) FollowUser(ctx context.Context, followerId string, followedId string) error {
	// --- Poslovna Logika ---
	// 1. Provera da li korisnik pokušava da zaprati sam sebe.
	if followerId == followedId {
		return ErrCannotFollowSelf
	}

	// Ovde bi mogle doći i druge provere, npr:
	// - Da li je nalog korisnika 'followedId' privatan?
	// - Da li je 'followerId' blokiran od strane 'followedId'?
	// Za sada, imamo samo osnovnu proveru.

	// Pozivamo repository da izvrši operaciju u bazi.
	return s.repo.FollowUser(ctx, followerId, followedId)
}

package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	repo      Repository
	jwtSecret interface{}
}

func NewService(r Repository, secret interface{}) *Service {
	return &Service{
		repo:      r,
		jwtSecret: secret,
	}
}

func (s *Service) MakeAuth(ctx context.Context, username string, password string) (string, error) {
	userID, err := s.repo.Auth(ctx, username, password)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"id":   userID,
		"pass": password,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := t.SignedString(s.jwtSecret)
	return accessToken, err
}

func (s *Service) PatchAuth(ctx context.Context, username string, password string) (string, error) {
	return "missing impl", nil // TODO
}

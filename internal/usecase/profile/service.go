package profile

import (
	"backend_coursework/internal/entity"
	"context"
)

type Service struct {
	repo ProfileRepo
}

func NewService(repo ProfileRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, user *entity.User) (entity.PK, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *Service) GetUserWithPolls(ctx context.Context, userID entity.PK, limit int) (*entity.User, error) {
	return s.repo.GetUserWithPolls(ctx, userID, limit)
}

func (s *Service) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, userID entity.PK) error {
	return s.repo.DeleteUser(ctx, userID)
}

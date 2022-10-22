package profile

import (
	"backend_coursework/internal/entity"
	"context"
)

type Service struct {
	repository ProfileRepo
}

func NewService(repository ProfileRepo) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUser(ctx context.Context, user *entity.User) error {
	return s.repository.CreateUser(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	return s.repository.GetUser(ctx, userID)
}

func (s *Service) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.repository.UpdateUser(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, userID entity.PK) error {
	return s.repository.DeleteUser(ctx, userID)
}

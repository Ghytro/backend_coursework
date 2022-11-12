package profile

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/repository"
	"context"
)

type ProfileRepo interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.PK, error)
	GetUser(ctx context.Context, userID entity.PK) (*entity.User, error)
	GetUserWithPolls(ctx context.Context, userID entity.PK, limit int) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID entity.PK) error

	RunInTransaction(ctx context.Context, f func(*database.TX) error) error
	WithTX(*database.TX) *repository.UserRepository
}

type PollsRepo interface {
	GetPollsCreatedBy(ctx context.Context, userID entity.PK, limit, offset int) ([]*entity.Poll, error)

	RunInTransaction(ctx context.Context, f func(*database.TX) error) error
	WithTX(*database.TX) *repository.PollsRepo
}

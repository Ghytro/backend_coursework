package profile

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"context"
)

type ProfileRepo interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.PK, error)
	GetUser(ctx context.Context, userID entity.PK) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID entity.PK) error

	RunInTransaction(ctx context.Context, f func(*database.TX) error) error
}

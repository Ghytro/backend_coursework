package auth

import (
	"backend_coursework/internal/entity"
	"context"
)

type UseCase interface {
	MakeAuth(ctx context.Context, username string, password string) (token string, err error)
	PatchAuth(ctx context.Context, username string, password string) (token string, err error)
	Register(ctx context.Context, user *entity.User) (token string, err error)
}

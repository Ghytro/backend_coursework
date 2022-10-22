package auth

import (
	"backend_coursework/internal/entity"
	"context"
)

type Repository interface {
	Auth(ctx context.Context, username string, password string) (entity.PK, error)
}

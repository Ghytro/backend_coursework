package auth

import "context"

type UseCase interface {
	MakeAuth(ctx context.Context, username string, password string) (token string, err error)
	PatchAuth(ctx context.Context, username string, password string) (token string, err error)
}

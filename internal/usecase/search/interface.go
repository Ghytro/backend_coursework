package search

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/repository"
	"context"
)

type UserRepository interface {
	GetUserListSearch(ctx context.Context, filter *repository.UserSearchFilter) ([]*entity.User, error)

	RunInTransaction(ctx context.Context, fn func(tx *database.TX) error) error
	WithTX(tx *database.TX) *repository.UserRepository
}

type PollsRepository interface {
	GetPollListSearch(ctx context.Context, filter *repository.PollSearchFilter) ([]*entity.Poll, error)

	RunInTransaction(ctx context.Context, fn func(tx *database.TX) error) error
	WithTX(tx *database.TX) *repository.PollsRepo
}

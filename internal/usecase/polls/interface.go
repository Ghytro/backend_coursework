package polls

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/repository"
	"context"
)

type Repository interface {
	CreatePoll(ctx context.Context, poll *entity.Poll) error
	GetPoll(ctx context.Context, id entity.PK) (*entity.Poll, error)
	GetPollWithVotes(ctx context.Context, id entity.PK) (*entity.Poll, error)
	GetPollWithVotesAmount(ctx context.Context, id entity.PK) (*entity.Poll, error)

	WithTX(*database.TX) *repository.PollsRepo
	RunInTransaction(context.Context, func(*database.TX) error) error
}

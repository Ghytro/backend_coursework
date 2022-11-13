package polls

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/repository"
	"context"
)

type Reader interface {
	GetPoll(ctx context.Context, id entity.PK) (*entity.Poll, error)
	GetPollCreator(ctx context.Context, id entity.PK) (*entity.User, error)
	GetVotesAmount(ctx context.Context, id entity.PK) ([]*entity.PollOption, error)
	GetUserPollVotes(ctx context.Context, userID entity.PK, pollID entity.PK) ([]*entity.Vote, error)
	GetPollListSearch(ctx context.Context, filter *repository.PollSearchFilter) ([]*entity.Poll, error)
	GetVoteListSearch(ctx context.Context, filter *repository.VoteSearchFilter) ([]*entity.Vote, error)
}

type Writer interface {
	CreatePoll(ctx context.Context, poll *entity.Poll) error
	Vote(ctx context.Context, userID entity.PK, pollID entity.PK, optIdxs ...int) error
	Unvote(ctx context.Context, userID entity.PK, pollID entity.PK) error
}

type Repository interface {
	Reader
	Writer

	WithTX(*database.TX) *repository.PollsRepo
	RunInTransaction(context.Context, func(*database.TX) error) error
}

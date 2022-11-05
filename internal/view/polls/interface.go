package polls

import (
	"backend_coursework/internal/entity"
	"context"
)

type UseCase interface {
	CreatePoll(ctx context.Context, creatorID entity.PK, model *NewPollRequest) (*entity.Poll, error)
	GetPollWithVotesAmount(ctx context.Context, id entity.PK, userID entity.PK) (*entity.Poll, []*entity.Vote, error)
	Vote(ctx context.Context, userID entity.PK, pollID entity.PK, optIdxs ...int) error
}

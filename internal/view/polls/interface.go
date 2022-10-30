package polls

import (
	"backend_coursework/internal/entity"
	"context"
)

type UseCase interface {
	CreatePoll(ctx context.Context, creatorID entity.PK, model *NewPollRequest) (*entity.Poll, error)
	GetPollWithVotesAmount(ctx context.Context, id entity.PK) (*entity.Poll, error)
}

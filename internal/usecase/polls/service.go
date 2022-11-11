package polls

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/validation"
	"backend_coursework/internal/view/polls"
	"context"

	"github.com/go-pg/pg/v10"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreatePoll(ctx context.Context, creatorID entity.PK, model *polls.NewPollRequest) (*entity.Poll, error) {
	if err := validation.ValidateCreatedPoll(model); err != nil {
		return nil, err
	}
	p := &entity.Poll{
		CreatorID:      creatorID,
		Topic:          model.Topic,
		IsAnonymous:    model.IsAnonymous == "on",
		MultipleChoice: model.MultipleChoice == "on",
		RevoteAbility:  model.CantRevote != "on",
	}
	for i, o := range model.Options {
		p.Options = append(p.Options, &entity.PollOption{
			Index:  i + 1,
			Option: o,
		})
	}
	err := s.repo.CreatePoll(ctx, p)
	return p, err
}

func (s *Service) GetPollWithVotesAmount(ctx context.Context, id entity.PK, userID entity.PK) (*entity.Poll, []*entity.Vote, error) {
	var (
		p         *entity.Poll
		userVotes []*entity.Vote
	)
	err := s.repo.RunInTransaction(ctx, func(tx *database.TX) error {
		repo := s.repo.WithTX(tx)
		var err error
		p, err = repo.GetPoll(ctx, id)
		if err != nil {
			return err
		}
		opts, err := repo.GetVotesAmount(ctx, id)
		if err != nil {
			return err
		}
		p.Options = opts
		creator, err := repo.GetPollCreator(ctx, id)
		if err != nil {
			return err
		}
		p.Creator = creator
		userVotes, err = repo.GetUserPollVotes(ctx, userID, id)
		if err == pg.ErrNoRows {
			err = nil
		}
		return err
	})
	return p, userVotes, err
}

func (s *Service) Vote(ctx context.Context, userID entity.PK, pollID entity.PK, optIdxs ...int) error {
	if len(optIdxs) == 0 {
		return nil
	}
	return s.repo.RunInTransaction(ctx, func(tx *database.TX) error {
		repo := s.repo.WithTX(tx)
		return repo.Vote(ctx, userID, pollID, optIdxs...)
	})
}

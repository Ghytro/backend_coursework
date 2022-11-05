package polls

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/view/polls"
	"context"
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/samber/lo"
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
	if model.Topic == "" {
		return nil, errors.New("тема опроса пуста")
	}
	if len(model.Options) == 0 {
		return nil, errors.New("у опроса нет вариантов ответа")
	}
	fmt.Println(*model)
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

		poll, err := repo.GetPoll(ctx, pollID)
		if err != nil {
			if err == pg.ErrNoRows {
				return fmt.Errorf("не найден опрос с id %d", pollID)
			}
			return err
		}

		if len(optIdxs) > 1 && !poll.MultipleChoice {
			return errors.New("на опросе не разрешен множественный выбор")
		}

		for _, i := range optIdxs {
			if !lo.ContainsBy(poll.Options, func(o *entity.PollOption) bool {
				return o.Index == i
			}) {
				return errors.New("выбрана опция, не присутствующая в опросе")
			}
		}

		return repo.Vote(ctx, userID, pollID, optIdxs...)
	})
}
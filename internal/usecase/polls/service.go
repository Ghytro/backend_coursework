package polls

import (
	"backend_coursework/internal/entity"
	"backend_coursework/internal/view/polls"
	"context"
	"errors"
	"fmt"
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

func (s *Service) GetPollWithVotesAmount(ctx context.Context, id entity.PK) (*entity.Poll, error) {
	return s.repo.GetPollWithVotesAmount(ctx, id)
}

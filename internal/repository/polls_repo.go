package repository

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/samber/lo"
)

type PollsRepo struct {
	db DBI
}

func NewPollsRepo(db DBI) *PollsRepo {
	return &PollsRepo{
		db: db,
	}
}

func (r *PollsRepo) WithTX(tx *database.TX) *PollsRepo {
	return NewPollsRepo(tx)
}

func (r *PollsRepo) RunInTransaction(ctx context.Context, fn func(*database.TX) error) error {
	return r.db.RunInTransaction(ctx, fn)
}

func (r *PollsRepo) CreatePoll(ctx context.Context, poll *entity.Poll) error {
	return r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		var u entity.User
		if err := tx.ModelContext(ctx, &u).Where("id = ? AND deleted_at IS NULL", poll.CreatorID).Select(); err != nil {
			return err
		}
		if _, err := tx.ModelContext(ctx, poll).Returning("*").Insert(); err != nil {
			return err
		}
		for i := range poll.Options {
			poll.Options[i].PollID = poll.ID
			poll.Options[i].Index = i + 1
		}
		_, err := tx.ModelContext(ctx, &poll.Options).Insert()
		return err
	})
}

func (r *PollsRepo) GetPoll(ctx context.Context, id entity.PK) (*entity.Poll, error) {
	resp := &entity.Poll{}
	resp.ID = id
	return resp, r.db.ModelContext(ctx, resp).WherePK().Relation("Options").Relation("Creator").Select()
}

func (r *PollsRepo) GetPollWithVotesAmount(ctx context.Context, id entity.PK) (*entity.Poll, error) {
	var (
		resp  entity.Poll
		model []struct {
			tableName struct{} `pg:"votes"`

			OptID entity.PK `pg:"option_id"`
			Count int       `pg:"cnt"`
		}
	)
	resp.ID = id
	err := r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		if err := tx.ModelContext(ctx, &resp).
			WherePK().
			Relation("Options").
			Select(); err != nil {
			return err
		}
		var u entity.User
		u.ID = resp.CreatorID
		if err := tx.ModelContext(ctx, &u).WherePK().Select(); err != nil {
			return err
		}
		resp.Creator = &u
		opts := lo.Map(resp.Options, func(opt *entity.PollOption, idx int) entity.PK {
			return opt.ID
		})
		return tx.ModelContext(ctx, &model).
			Where("option_id IN (?)", pg.In(opts)).
			Group("option_id").
			ColumnExpr("option_id, count(*) AS cnt").
			Select()
	})
	if err != nil {
		return nil, err
	}
	m := make(map[entity.PK]int)
	for _, o := range model {
		m[o.OptID] = o.Count
	}
	for i := range resp.Options {
		resp.Options[i].VotesAmount = m[resp.Options[i].ID]
	}
	return &resp, nil
}

func (r *PollsRepo) GetPollWithVotes(ctx context.Context, id entity.PK) (*entity.Poll, error) {
	var resp entity.Poll
	resp.ID = id
	err := r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		if err := tx.ModelContext(ctx, &resp).WherePK().Relation("Options").Relation("Creator").Select(); err != nil {
			return err
		}
		return tx.ModelContext(ctx, &resp.Options).Where("poll_id = ?", id).Relation("Votes").Select()
	})
	return &resp, err
}

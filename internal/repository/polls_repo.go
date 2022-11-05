package repository

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"context"
	"fmt"

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
	return resp, r.db.ModelContext(ctx, resp).WherePK().Relation("Options").Select()
}

func (r *PollsRepo) GetPollCreator(ctx context.Context, id entity.PK) (*entity.User, error) {
	var p entity.Poll
	p.ID = id
	return p.Creator, r.db.ModelContext(ctx, &p).WherePK().Relation("Creator").Select()
}

func (r *PollsRepo) GetVotesAmount(ctx context.Context, id entity.PK) ([]*entity.PollOption, error) {
	var (
		options []*entity.PollOption
		model   []struct {
			tableName struct{} `pg:"votes"`

			OptID entity.PK `pg:"option_id"`
			Count int       `pg:"cnt"`
		}
	)
	r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		if err := tx.ModelContext(ctx, &options).Where("poll_id = ?", id).Select(); err != nil {
			return err
		}
		optIds := lo.Map(options, func(o *entity.PollOption, idx int) entity.PK {
			return o.ID
		})

		return tx.ModelContext(ctx, &model).
			Where("option_id IN (?)", pg.In(optIds)).
			Group("option_id").
			ColumnExpr("option_id, count(*) AS cnt").
			Select()
	})
	for _, o := range model {
		opt, _ := lo.Find(options, func(el *entity.PollOption) bool {
			return el.ID == o.OptID
		})
		opt.VotesAmount = o.Count
	}
	return options, nil
}

func (r *PollsRepo) UserVoted(ctx context.Context, userID entity.PK, pollID entity.PK) (bool, error) {
	var vote entity.Vote
	if err := r.db.ModelContext(ctx, &vote).Where("user_id = ? AND poll_id = ?", userID, pollID).Select(); err != nil {
		if err == pg.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *PollsRepo) GetUserPollVotes(ctx context.Context, userID entity.PK, pollID entity.PK) ([]*entity.Vote, error) {
	var result []*entity.Vote
	return result, r.db.ModelContext(ctx, &result).Where("user_id = ? AND poll_id = ?", userID, pollID).Select()
}

func (r *PollsRepo) Vote(ctx context.Context, userID entity.PK, pollID entity.PK, optionIdxs ...int) error {
	return r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		var u entity.User
		u.ID = userID
		if err := tx.ModelContext(ctx, &u).WherePK().Select(); err != nil {
			return fmt.Errorf("не найден пользователь с id %d", u.ID)
		}
		var p entity.Poll
		p.ID = pollID
		if err := tx.ModelContext(ctx, &p).WherePK().Relation("Options").Select(); err != nil {
			return fmt.Errorf("не найден опрос с id %d", p.ID)
		}
		p.Options = lo.Filter(p.Options, func(opt *entity.PollOption, _ int) bool {
			return lo.Contains(optionIdxs, opt.Index)
		})
		v := lo.Map(p.Options, func(opt *entity.PollOption, _ int) *entity.Vote {
			return &entity.Vote{
				UserID:   userID,
				PollID:   pollID,
				OptionID: opt.ID,
			}
		})
		_, err := tx.ModelContext(ctx, &v).Insert()
		return err
	})
}

func (r *PollsRepo) Unvote(ctx context.Context, userID entity.PK, pollID entity.PK) error {
	return r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		var p entity.Poll
		p.ID = pollID
		if err := tx.ModelContext(ctx, &p).WherePK().Relation("Options").Select(); err != nil {
			return err
		}
		var v []*entity.Vote
		_, err := tx.ModelContext(ctx, &v).Where("poll_id = ? AND user_id = ?", pollID, userID).Delete()
		return err
	})
}

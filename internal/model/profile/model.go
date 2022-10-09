package profile

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
)

type Model struct {
	db database.DBI
}

func NewProfileModel(db database.DBI) *Model {
	return &Model{
		db: db,
	}
}

func (m *Model) WithTX(tx database.DBI) *Model {
	return NewProfileModel(tx)
}

func (m *Model) CreateUser(ctx context.Context, user *entity.User) error {
	return m.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		var u entity.User
		if err := tx.ModelContext(ctx, &u).Where("username = ? AND deleted_at IS NULL", user.Username).Select(); err != nil {
			if err == pg.ErrNoRows {
				_, err := tx.ModelContext(ctx, user).Returning("*").Insert()
				return err
			}
		}
		return errors.New("пользователь с таким именем уже существует")
	})
}

func (m *Model) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	var u entity.User
	if err := m.db.ModelContext(ctx, &u).Where("id = ? AND deleted_at IS NULL ", userID).
		Relation("Polls").Relation("Votes").Select(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *Model) UpdateUser(ctx context.Context, user *entity.User) error {
	_, err := m.db.ModelContext(ctx, user).Where("id = ? AND deleted_at IS NULL", user.ID).Update()
	return err
}

func (m *Model) DeleteUser(ctx context.Context, userID entity.PK) error {
	_, err := m.db.ModelContext(ctx, (*entity.User)(nil)).Set("deleted_at = NOW()").Where("id = ?", userID).Update()
	return err
}

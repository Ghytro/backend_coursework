package repository

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
)

type UserRepository struct {
	db DBI
}

func NewUserRepo(db DBI) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (m *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return m.db.RunInTransaction(ctx, func(tx *database.TX) error {
		var u entity.User
		if err := tx.ModelContext(ctx, &u).Where("username = ? AND deleted_at IS NULL", user.Username).Select(); err != nil {
			if err == pg.ErrNoRows {
				_, err := tx.ModelContext(ctx, user).Value("password", "crypt(?password, gen_salt('bf'))").Insert()
				return err
			}
		}
		return errors.New("пользователь с таким именем уже существует")
	})
}

func (m *UserRepository) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	var u entity.User
	if err := m.db.ModelContext(ctx, &u).Where("id = ? AND deleted_at IS NULL ", userID).
		Relation("Polls").Relation("Votes").Select(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	_, err := m.db.ModelContext(ctx, user).Where("id = ? AND deleted_at IS NULL", user.ID).Update()
	return err
}

func (m *UserRepository) DeleteUser(ctx context.Context, userID entity.PK) error {
	_, err := m.db.ModelContext(ctx, (*entity.User)(nil)).Set("deleted_at = NOW()").Where("id = ?", userID).Update()
	return err
}

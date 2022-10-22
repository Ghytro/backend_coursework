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

func (m *UserRepository) RunInTransaction(ctx context.Context, f func(tx *database.TX) error) error {
	return m.db.RunInTransaction(ctx, f)
}

func (m *UserRepository) WithTX(tx *database.TX) *UserRepository {
	return NewUserRepo(tx)
}

func (m *UserRepository) CreateUser(ctx context.Context, user *entity.User) (entity.PK, error) {
	var u entity.User
	err := m.db.RunInTransaction(ctx, func(tx *database.TX) error {
		if err := tx.ModelContext(ctx, &u).Where("username = ?", user.Username).Select(); err != nil {
			if err == pg.ErrNoRows {
				_, _err := tx.ModelContext(ctx, &u).
					Value("username", "?", user.Username).
					Value("password", "crypt(?, gen_salt('bf'))", user.Password).
					Value("first_name", "?", user.FirstName).
					Value("last_name", "?", user.LastName).
					Value("bio", "?", user.Bio).
					Value("avatar_url", "?", user.AvatarUrl).
					Value("country", "?", user.Country).
					Returning("*").Insert()
				return _err
			}
			return err
		}
		if u.DeletedAt.IsZero() {
			return errors.New("пользователь с таким именем уже существует")
		}
		_, err := tx.ModelContext(ctx, &u).Where("id = ?", u.ID).Set(
			`username = ?,
			first_name = ?,
			last_name = ?,
			password = crypt(?, password),
			bio = ?,
			avatar_url = ?,
			country = ?,
			created_at = NOW(),
			deleted_at = NULL`,
			user.Username,
			user.FirstName,
			user.LastName,
			user.Password,
			user.Bio,
			user.AvatarUrl,
			user.Country,
		).Update()
		return err
	})
	return u.ID, err
}

func (m *UserRepository) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	var u entity.User
	if err := m.db.ModelContext(ctx, &u).Where("id = ? AND deleted_at IS NULL", userID).
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

func (m *UserRepository) Auth(ctx context.Context, username string, password string) (entity.PK, error) {
	var u entity.User
	if err := m.db.ModelContext(ctx, &u).Where("username = ? AND password = crypt(?, password)", username, password).Select(); err != nil {
		return 0, err
	}
	return u.ID, nil
}

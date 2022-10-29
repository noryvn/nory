package user

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"nory/domain"
)

type UserRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewUserRepositoryPostgres(pool *pgxpool.Pool) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{pool}
}

func (urp *UserRepositoryPostgres) GetUserByUserId(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{
		UserId:    id,
		CreatedAt: time.Time{},
		Username:  "",
		Name:      "",
		Email:     "",
	}
	row := urp.pool.QueryRow(ctx, "SELECT username, name, email, created_at FROM app_user WHERE user_id = $1", id)
	err := row.Scan(
		&u.Username,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrUserNotExists
	}
	return u, err
}

func (urp *UserRepositoryPostgres) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	u := &domain.User{
		UserId:    "",
		CreatedAt: time.Time{},
		Username:  username,
		Name:      "",
		Email:     "",
	}
	row := urp.pool.QueryRow(ctx, "SELECT user_id, name, email, created_at FROM app_user WHERE username = $1", username)
	err := row.Scan(
		&u.UserId,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrUserNotExists
	}
	return u, err
}

func (urp *UserRepositoryPostgres) DeleteUser(ctx context.Context, id string) error {
	_, err := urp.pool.Exec(
		ctx,
		"DELETE FROM app_user WHERE user_id = $1",
		id,
	)
	return err
}

func (urp *UserRepositoryPostgres) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := urp.pool.Exec(
		ctx,
		"INSERT INTO app_user(user_id, username, name, email) VALUES($1, $2, $3, $4)",
		user.UserId,
		user.Username,
		user.Name,
		user.Email,
	)
	if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
		return domain.ErrUserAlreadyExists
	}
	return err
}

func (urp *UserRepositoryPostgres) UpdateUser(ctx context.Context, user *domain.User) error {
	u, err := urp.GetUserByUserId(ctx, user.UserId)
	if err != nil {
		return err
	}
	u.Update(user)
	_, err = urp.pool.Exec(
		ctx,
		`UPDATE app_user SET username = $1, name = $2 WHERE user_id = $3`,
		u.Username,
		u.Name,
		u.UserId,
	)
	if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
		return domain.ErrUserAlreadyExists
	}
	return err
}

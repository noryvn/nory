package user_test

import (
	"context"
	"os"
	"testing"

	"nory/domain"
	. "nory/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var (
	userFoo = uuid.NewString()
	userBar = uuid.NewString()
	userBaz = uuid.NewString()
	userQux = uuid.NewString()
)

func TestUserRepository(t *testing.T) {
	t.Parallel()
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Error(err)
	}
	repos := []Repository{
		{
			Name:           "postgres",
			UserRepository: NewUserRepositoryPostgres(pool),
			Skip:           os.Getenv("DATABASE_URL") == "",
		},
		{
			Name:           "memory",
			UserRepository: NewUserRepositoryMem(),
			Skip:           false,
		},
	}

	for _, repo := range repos {
		repo := repo
		t.Run(repo.Name, func(t *testing.T) {
			if repo.Skip {
				t.Skipf("skipping %s", repo.Name)
			}
			t.Parallel()
			t.Run("CreateUser", repo.testCreateUser)
			t.Run("GetUser", repo.testGetUser)
			t.Run("UpdateUser", repo.testUpdateUser)
			t.Run("DeleteUser", repo.testDeleteUser)
		})
	}
}

type Repository struct {
	Name           string
	UserRepository domain.UserRepository
	Skip           bool
}

func (r *Repository) testCreateUser(t *testing.T) {
	testCases := []struct {
		Name string
		User domain.User
		Err  error
	}{
		{"success create", domain.User{Username: "foo", Email: "foo@bel.ia", UserId: userFoo}, nil},
		{"success create", domain.User{Username: "bar", Email: "bar@bel.ia", UserId: userBar}, nil},
		{"success create", domain.User{Username: "baz", Email: "baz@bel.ia", UserId: userBaz}, nil},
		{"duplicate username", domain.User{Username: "foo", UserId: userQux}, domain.ErrUserAlreadyExists},
		{"duplicate id", domain.User{UserId: userFoo}, domain.ErrUserAlreadyExists},
		{"duplicate email", domain.User{UserId: userQux, Email: "foo@bel.ia"}, domain.ErrUserAlreadyExists},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			err := r.UserRepository.CreateUser(context.Background(), &tc.User)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if err == nil {
				u, err := r.UserRepository.GetUser(context.Background(), tc.User.UserId)
				assert.Nil(t, err)
				uu, err := r.UserRepository.GetUserWithUsername(context.Background(), tc.User.Username)
				assert.Nil(t, err)
				assert.Equal(t, u, uu)
			}
		})
	}
}

func (r *Repository) testGetUser(t *testing.T) {
	testCases := []struct {
		Name string
		Id   string
		Err  error
	}{
		{"success", userFoo, nil},
		{"failed", userQux, domain.ErrUserNotExists},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			u, err := r.UserRepository.GetUser(context.Background(), tc.Id)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if tc.Err == nil && err == nil {
				assert.Equal(t, tc.Id, u.UserId, "missmatch user id")
				uu, err := r.UserRepository.GetUserWithUsername(context.Background(), u.Username)
				assert.Nil(t, err)
				assert.Equal(t, u, uu)
			}
		})
	}
}

func (r *Repository) testUpdateUser(t *testing.T) {
	testCases := []struct {
		Name string
		User domain.User
		Err  error
	}{
		{"success", domain.User{UserId: userFoo, Username: "foo-bar", Email: "a@bel.ia"}, nil},
		{"duplicate username", domain.User{UserId: userBar, Username: "foo-bar"}, domain.ErrUserAlreadyExists},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			prev, err := r.UserRepository.GetUser(context.Background(), tc.User.UserId)
			assert.Equal(t, nil, err, "unexpected error received")

			err = r.UserRepository.UpdateUser(context.Background(), &tc.User)
			assert.Equal(t, tc.Err, err, "missmatch error")

			curr, err := r.UserRepository.GetUser(context.Background(), tc.User.UserId)
			assert.Equal(t, nil, err, "unexpected error received")

			if tc.Err == nil && err == nil {
				assert.Equal(t, prev.UserId, curr.UserId, "update should not change user id")
				assert.Equal(t, prev.CreatedAt, curr.CreatedAt, "update should not change created at")
				assert.Equal(t, tc.User.Username, curr.Username)
			}
		})
	}
}

func (r *Repository) testDeleteUser(t *testing.T) {
	testCases := []struct {
		Name string
		Id   string
		Err  error
	}{
		{"delete existing user", userFoo, nil},
		{"delete existing user", userBar, nil},
		{"delete existing user", userBaz, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			user, err := r.UserRepository.GetUser(context.Background(), tc.Id)
			assert.Nil(t, err)
			err = r.UserRepository.DeleteUser(context.Background(), tc.Id)
			assert.Equal(t, nil, err, "missmatch error")
			_, err = r.UserRepository.GetUser(context.Background(), tc.Id)
			assert.Equal(t, domain.ErrUserNotExists, err)
			_, err = r.UserRepository.GetUserWithUsername(context.Background(), user.Username)
			assert.Equal(t, domain.ErrUserNotExists, err)
			err = r.UserRepository.DeleteUser(context.Background(), tc.Id)
			assert.Equal(t, nil, err, "missmatch error")
		})
	}
}

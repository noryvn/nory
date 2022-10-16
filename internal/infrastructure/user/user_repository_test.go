package user_test

import (
	"context"
	"testing"

	"nory/internal/domain"
	. "nory/internal/infrastructure/user"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
	repos := []Repository{
		{
			Name: "memory",
			R:    NewUserRepositoryMem(),
		},
	}

	for _, repo := range repos {
		repo := repo
		t.Run(repo.Name, func(t *testing.T) {
			t.Parallel()
			t.Run("create", repo.testCreate)
			t.Run("get", repo.testGet)
			t.Run("update", repo.testUpdate)
			t.Run("delete", repo.testDelete)
		})
	}
}

type Repository struct {
	Name string
	R    domain.UserRepository
}

func (r Repository) testCreate(t *testing.T) {
	testCases := []struct {
		Name string
		User domain.User
		Err  error
	}{
		{"success create", domain.User{Username: "foo", UserId: "foo"}, nil},
		{"success create", domain.User{Username: "bar", UserId: "bar"}, nil},
		{"success create", domain.User{Username: "baz", UserId: "baz"}, nil},
		{"duplicate username", domain.User{Username: "foo"}, domain.ErrDuplicateUser},
		{"duplicate id", domain.User{UserId: "foo"}, domain.ErrUserExists},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			err := r.R.CreateUser(context.Background(), &tc.User)
			assert.Equal(t, tc.Err, err, "missmatch error")
		})
	}
}

func (r Repository) testGet(t *testing.T) {
	testCases := []struct {
		Name string
		Id   string
		Err  error
	}{
		{"success", "foo", nil},
		{"success", "bar", nil},
		{"failed", "qux", domain.ErrUserNotFound},
		{"failed", "hai", domain.ErrUserNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			u, err := r.R.GetUser(context.Background(), tc.Id)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if tc.Err == nil {
				assert.Equal(t, tc.Id, u.UserId, "missmatch user id")
			}
		})
	}
}

func (r Repository) testUpdate(t *testing.T) {
	testCases := []struct {
		Name string
		User domain.User
		Err  error
	}{
		{"success", domain.User{UserId: "foo", Username: "foo-bar"}, nil},
		{"duplicate username", domain.User{UserId: "bar", Username: "foo-bar"}, domain.ErrDuplicateUser},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			prev, err := r.R.GetUser(context.Background(), tc.User.UserId)
			assert.Equal(t, nil, err, "unexpected error received")

			err = r.R.UpdateUser(context.Background(), &tc.User)
			assert.Equal(t, tc.Err, err, "missmatch error")

			curr, err := r.R.GetUser(context.Background(), tc.User.UserId)
			assert.Equal(t, nil, err, "unexpected error received")

			if tc.Err == nil {
				assert.Equal(t, prev.UserId, curr.UserId, "update should not change user id")
				assert.Equal(t, prev.CreatedAt, curr.CreatedAt, "update should not change created at")
				assert.Equal(t, tc.User.Username, curr.Username)
			}
		})
	}
}

func (r Repository) testDelete(t *testing.T) {
	testCases := []struct {
		Name string
		Id   string
		Err  error
	}{
		{"delete existing user", "foo", nil},
		{"delete unexists", "foo", nil},
		{"delete unexists", "foo-bar-baz", nil},
		{"delete existing user", "bar", nil},
		{"delete existing user", "baz", nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			err := r.R.DeleteUser(context.Background(), tc.Id)
			assert.Equal(t, tc.Err, err, "missmatch error")
		})
	}
}

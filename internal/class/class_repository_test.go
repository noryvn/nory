package class_test

import (
	"context"
	"os"
	"testing"

	"nory/domain"
	. "nory/internal/class"
	"nory/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

const (
	userFoo = "09229265-d53c-44fd-a791-cb686b3e61d6"
	userBar = "eff59eb4-ef31-4491-bc7e-77026f4cb5a8"
	userBaz = "ff4a1fc5-b7bf-40e4-be55-0eb9dfb886f7"
	userQux = "cde9b03c-311d-443d-901c-45473f453305"
)

func TestClassRepository(t *testing.T) {
	t.Parallel()
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Error(err)
	}
	repos := []Repository{
		{
			Name: "memory",
			R:    NewClassRepositoryMem(),
		},
		{
			Name: "postgres",
			R:    NewClassRepositoryPostgres(pool),
			Skip: os.Getenv("DATABASE_URL") == "",
		},
	}

	for _, r := range repos {
		r := r
		t.Run(r.Name, func(t *testing.T) {
			if r.Skip {
				t.Skipf("skipping %s", r.Name)
			}
			if r.Name == "postgres" {
				userRepo := user.NewUserRepositoryPostgres(pool)
				assert.Nil(t, userRepo.CreateUser(context.Background(), &domain.User{UserId: userFoo, Email: userFoo, Username: "userFoo" }))
				assert.Nil(t, userRepo.CreateUser(context.Background(), &domain.User{UserId: userBar, Email: userBar, Username: "userBar" }))
				assert.Nil(t, userRepo.CreateUser(context.Background(), &domain.User{UserId: userBaz, Email: userBaz, Username: "userBaz" }))
				t.Cleanup(func() {
					userRepo.DeleteUser(context.Background(), userFoo)
					userRepo.DeleteUser(context.Background(), userBar)
					userRepo.DeleteUser(context.Background(), userBaz)
				})
			}
			t.Parallel()
			t.Run("create", r.testCreate)
			t.Run("get by class id", r.testGet)
			t.Run("get by owner id", r.testGetByOwnerId)
			t.Run("update class", r.testUpdate)
			t.Run("delete", r.testDelete)
		})
	}
}

type Repository struct {
	Name string
	R    domain.ClassRepository
	Skip bool
	Classes []domain.Class
}

func (r *Repository) testCreate(t *testing.T) {
	testCases := []struct {
		Name  string
		Class domain.Class
		Err   error
	}{
		{"success", domain.Class{ClassId: "foo", OwnerId: userFoo}, nil},
		{"success", domain.Class{ClassId: "bar", OwnerId: userFoo}, nil},
		{"success", domain.Class{ClassId: "baz", OwnerId: userBar}, nil},
		{"success", domain.Class{ClassId: "baz", OwnerId: userBaz}, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			class := tc.Class
			err := r.R.CreateClass(context.Background(), &class)
			assert.Equal(t, tc.Err, err, "missmatch err")
			if err == nil {
				assert.NotEqual(t, tc.Class.ClassId, class.ClassId, "CreateClass must update (Class).ClassId to generated id")
				r.Classes = append(r.Classes, class)
			}
		})
	}
}

func (r *Repository) testDelete(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		Err     error
	}{
		{"existing class", r.Classes[0].ClassId, nil},
		{"existing class", r.Classes[1].ClassId, nil},
		{"existing class", r.Classes[2].ClassId, nil},
		{"unexisting class", "baz", nil},
		{"unexisting class", "baz", nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			err := r.R.DeleteClass(context.Background(), tc.ClassId)
			assert.Equal(t, tc.Err, err, "missmatch err")
		})
	}
}

func (r *Repository) testGet(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		Err     error
	}{
		{"existing class", r.Classes[0].ClassId, nil},
		{"existing class", r.Classes[1].ClassId, nil},
		{"existing class", r.Classes[2].ClassId, nil},
		{"unexisting class", "foo-bar", domain.ErrClassNotFound},
		{"unexisting class", "foo-baz", domain.ErrClassNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			c, err := r.R.GetClass(context.Background(), tc.ClassId)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if err == nil {
				assert.Equal(t, tc.ClassId, c.ClassId, "unexpected class id")
			}
		})
	}
}

func (r *Repository) testGetByOwnerId(t *testing.T) {
	testCases := []struct {
		Name    string
		OwnerId string
		Len     int
		Err     error
	}{
		{"exists", userFoo, 2, nil},
		{"exists", userBar, 1, nil},
		{"unexists", userQux, 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			classes, err := r.R.GetClassesByOwnerId(context.Background(), tc.OwnerId)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if err == nil {
				assert.Equal(t, tc.Len, len(classes), "unexpected class received")
				for _, c := range classes {
					assert.Equal(t, tc.OwnerId, c.OwnerId, "unexpected class owner")
				}
			}
		})
	}
}

func (r *Repository) testUpdate(t *testing.T) {
	testCases := []struct {
		Name  string
		Class domain.Class
		Err   error
	}{
		{"success", domain.Class{ClassId: r.Classes[0].ClassId, Description: "foo"}, nil},
		{"not found", domain.Class{ClassId: "anu", Description: "foo"}, domain.ErrClassNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			prev, err := r.R.GetClass(context.Background(), tc.Class.ClassId)
			assert.Equal(t, tc.Err, err, "unexpected error")

			err = r.R.UpdateClass(context.Background(), &tc.Class)
			assert.Equal(t, tc.Err, err, "missmatch error")

			curr, err := r.R.GetClass(context.Background(), tc.Class.ClassId)
			assert.Equal(t, tc.Err, err, "unexpected error")

			if err == nil {
				assert.Equal(t, prev.ClassId, curr.ClassId, "should not update class id")
				assert.Equal(t, prev.OwnerId, curr.OwnerId, "should not update owner id")
				assert.Equal(t, tc.Class.Description, curr.Description, "should able update Description")
			}
		})
	}
}

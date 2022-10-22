package class_test

import (
	"context"
	"os"
	"testing"

	"nory/domain"
	. "nory/internal/class"
	"nory/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestClassRepository(t *testing.T) {
	t.Parallel()
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Error(err)
	}
	repos := []Repository{
		{
			Name:            "memory",
			ClassRepository: NewClassRepositoryMem(),
			UserRepository:  user.NewUserRepositoryMem(),
		},
		{
			Name:            "postgres",
			ClassRepository: NewClassRepositoryPostgres(pool),
			UserRepository:  user.NewUserRepositoryPostgres(pool),
			Skip:            os.Getenv("DATABASE_URL") == "",
		},
	}

	for _, r := range repos {
		r := r
		t.Run(r.Name, func(t *testing.T) {
			r.t = t
			if r.Skip {
				t.Skipf("skipping %s", r.Name)
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
	Name            string
	ClassRepository domain.ClassRepository
	UserRepository  domain.UserRepository
	Skip            bool

	t       *testing.T
	classes []domain.Class
	users   map[string]string
}

func (r *Repository) getUser(name string) string {
	if r.users == nil {
		r.users = make(map[string]string)
	}
	if id, ok := r.users[name]; ok {
		return id
	}
	id := uuid.NewString()
	r.users[name] = id
	user := &domain.User{
		UserId:   id,
		Email:    xid.New().String(),
		Username: xid.New().String(),
	}
	err := r.UserRepository.CreateUser(context.Background(), user)
	assert.Nil(r.t, err, "failed to create user")
	r.t.Cleanup(func() {
		r.UserRepository.DeleteUser(context.Background(), id)
	})
	return id
}

func (r *Repository) testCreate(t *testing.T) {
	foo := r.getUser("foo")
	bar := r.getUser("bar")
	baz := r.getUser("baz")
	testCases := []struct {
		Name  string
		Class domain.Class
		Err   error
	}{
		{"success", domain.Class{ClassId: "foo", OwnerId: foo}, nil},
		{"success", domain.Class{ClassId: "bar", OwnerId: foo}, nil},
		{"success", domain.Class{ClassId: "baz", OwnerId: bar}, nil},
		{"success", domain.Class{ClassId: "baz", OwnerId: baz}, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			class := tc.Class
			err := r.ClassRepository.CreateClass(context.Background(), &class)
			assert.Equal(t, tc.Err, err, "missmatch err")
			if err == nil {
				assert.NotEqual(t, tc.Class.ClassId, class.ClassId, "CreateClass must update (Class).ClassId to generated id")
				r.classes = append(r.classes, class)
				t.Log(class)
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
		{"existing class", r.classes[0].ClassId, nil},
		{"existing class", r.classes[1].ClassId, nil},
		{"existing class", r.classes[2].ClassId, nil},
		{"unexisting class", "baz", nil},
		{"unexisting class", "baz", nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			err := r.ClassRepository.DeleteClass(context.Background(), tc.ClassId)
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
		{"existing class", r.classes[0].ClassId, nil},
		{"existing class", r.classes[1].ClassId, nil},
		{"existing class", r.classes[2].ClassId, nil},
		{"unexisting class", "foo-bar", domain.ErrClassNotExists},
		{"unexisting class", "foo-baz", domain.ErrClassNotExists},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			c, err := r.ClassRepository.GetClass(context.Background(), tc.ClassId)
			t.Log(c, tc.ClassId, err)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if err == nil {
				assert.NotNil(t, c)
				assert.Equal(t, tc.ClassId, c.ClassId, "unexpected class id")
			}
		})
	}
}

func (r *Repository) testGetByOwnerId(t *testing.T) {
	foo := r.getUser("foo")
	bar := r.getUser("bar")
	qux := r.getUser("qux")
	testCases := []struct {
		Name    string
		OwnerId string
		Len     int
		Err     error
	}{
		{"exists", foo, 2, nil},
		{"exists", bar, 1, nil},
		{"unexists", qux, 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			classes, err := r.ClassRepository.GetClassesByOwnerId(context.Background(), tc.OwnerId)
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
		{"success", domain.Class{ClassId: r.classes[0].ClassId, Description: "foo"}, nil},
		{"not found", domain.Class{ClassId: "anu", Description: "foo"}, domain.ErrClassNotExists},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			prev, err := r.ClassRepository.GetClass(context.Background(), tc.Class.ClassId)
			assert.Equal(t, tc.Err, err, "unexpected error")

			err = r.ClassRepository.UpdateClass(context.Background(), &tc.Class)
			assert.Equal(t, tc.Err, err, "missmatch error")

			curr, err := r.ClassRepository.GetClass(context.Background(), tc.Class.ClassId)
			assert.Equal(t, tc.Err, err, "unexpected error")

			if err == nil {
				assert.Equal(t, prev.ClassId, curr.ClassId, "should not update class id")
				assert.Equal(t, prev.OwnerId, curr.OwnerId, "should not update owner id")
				assert.Equal(t, tc.Class.Description, curr.Description, "should able update Description")
			}
		})
	}
}

package classmember_test

import (
	"context"
	"os"
	"testing"
	"time"

	"nory/domain"
	"nory/internal/class"
	. "nory/internal/class_member"
	"nory/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestClassMemberRepository(t *testing.T) {
	t.Parallel()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	assert.Nil(t, err)

	for _, repo := range []Repository{
		{
			Name: "memory",
			Repo: NewClassMemberRepositoryMem(), Skip: false},
		{
			Name: "postgres",
			Repo: NewClassMemberRepositoryPostgres(pool),
			ClassRepo: class.NewClassRepositoryPostgres(pool),
			UserRepo: user.NewUserRepositoryPostgres(pool),
			Skip: os.Getenv("DATABASE_URL") == "",
		},
	} {
		repo := repo
		t.Run(repo.Name, func(t *testing.T) {
			t.Parallel()
			classFoo := xid.New().String()
			classBar := xid.New().String()
			userFoo := uuid.NewString()
			userBar := uuid.NewString()

			if repo.ClassRepo != nil {
				owner := uuid.NewString()
				err := repo.UserRepo.CreateUser(context.Background(), &domain.User{
					UserId:    userFoo,
					Username:  xid.New().String(),
					Email:     xid.New().String(),
				})
				assert.Nil(t, err)

				err = repo.UserRepo.CreateUser(context.Background(), &domain.User{
					UserId:    userBar,
					Username:  xid.New().String(),
					Email:     xid.New().String(),
				})
				assert.Nil(t, err)

				err = repo.UserRepo.CreateUser(context.Background(), &domain.User{
					UserId:    owner,
					Username:  xid.New().String(),
					Email:     xid.New().String(),
				})
				assert.Nil(t, err)

				c := &domain.Class{
					ClassId:     "",
					OwnerId:     owner,
				}

				err = repo.ClassRepo.CreateClass(context.Background(), c)
				assert.Nil(t, err)
				classFoo = c.ClassId

				err = repo.ClassRepo.CreateClass(context.Background(), c)
				assert.Nil(t, err)
				classBar = c.ClassId
			}

			memA := domain.ClassMember{UserId: userFoo, ClassId: classBar, Level: "member"}
			memB := domain.ClassMember{UserId: userFoo, ClassId: classFoo, Level: "member"}
			memC := domain.ClassMember{UserId: userBar, ClassId: classFoo, Level: "member"}

			for _, m := range []domain.ClassMember{
				memA,
				memB,
				memC,
			} {
				m := m
				updated := m
				updated.Level = "admin"
				err := repo.Repo.CreateMember(context.Background(), &m)
				if assert.Nil(t, err) {
					err := repo.Repo.UpdateMember(context.Background(), &updated)
					assert.Nil(t, err)

					mem, err := repo.Repo.GetMember(context.Background(), &m)
					assert.Nil(t, err)
					mem.CreatedAt = time.Time{}
					assert.Equal(t, updated, *mem)

					t.Cleanup(func() {
						err := repo.Repo.DeleteMember(context.Background(), &m)
						assert.Nil(t, err)

						err = repo.Repo.DeleteMember(context.Background(), &m)
						assert.Nil(t, err)

						_, err = repo.Repo.GetMember(context.Background(), &m)
						assert.Equal(t, domain.ErrClassMemberNotExists, err)
					})
				}
			}

			t.Run("list members", func(t *testing.T) {
				for _, tc := range []struct {
					ClassId string
					Len     int
				}{
					{classFoo, 2},
					{classBar, 1},
					{xid.New().String(), 0},
				} {
					members, err := repo.Repo.ListMembers(context.Background(), tc.ClassId)
					assert.Nil(t, err)
					assert.Equal(t, tc.Len, len(members))
				}
			})

			t.Run("list joined", func(t *testing.T) {
				for _, tc := range []struct {
					UserId string
					Len    int
				}{
					{userFoo, 2},
					{userBar, 1},
					{xid.New().String(), 0},
				} {
					members, err := repo.Repo.ListJoined(context.Background(), tc.UserId)
					assert.Nil(t, err)
					assert.Equal(t, tc.Len, len(members))
				}
			})
		})
	}
}

type Repository struct {
	Name string
	Repo domain.ClassMemberRepository
	ClassRepo domain.ClassRepository
	UserRepo domain.UserRepository
	Skip bool
}

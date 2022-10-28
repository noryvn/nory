package classmember_test

import (
	"context"
	"testing"

	"nory/domain"
	. "nory/internal/class_member"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)


func TestClassMemberRepository(t *testing.T) {
	t.Parallel()

	for _, repo := range []Repository{
		{"memory", NewClassMemberRepositoryMem()},
	}{
		repo := repo
		t.Run(repo.Name, func(t *testing.T) {
			t.Parallel()
			classFoo := xid.New().String()
			classBar := xid.New().String()
			userFoo := xid.New().String()
			userBar := xid.New().String()

			memA := domain.ClassMember{UserId: userFoo, ClassId: classBar, Level: "member"}
			memB := domain.ClassMember{UserId: userFoo, ClassId: classFoo, Level: "member"}
			memC := domain.ClassMember{UserId: userBar, ClassId: classFoo, Level: "member"}

			for _, m := range []domain.ClassMember{
				memA,
				memB,
				memC,
			}{
				m := m
				updated := m
				updated.Level = "admin"
				err := repo.Repo.CreateMember(context.Background(), &m)
				if assert.Nil(t, err) {
					err := repo.Repo.UpdateMember(context.Background(), &updated)
					assert.Nil(t, err)

					mem, err := repo.Repo.GetMember(context.Background(), &m)
					assert.Nil(t, err)
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
				for _, tc := range []struct{
					ClassId string
					Len int
				}{
					{classFoo, 2},
					{classBar, 1},
					{xid.New().String(), 0},
				}{
					members, err := repo.Repo.ListMembers(context.Background(), tc.ClassId)
					assert.Nil(t, err)
					assert.Equal(t, tc.Len, len(members))
				}
			})

			t.Run("list joined", func(t *testing.T) {
				for _, tc := range []struct{
					UserId string
					Len int
				}{
					{userFoo, 2},
					{userBar, 1},
					{xid.New().String(), 0},
				}{
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
}

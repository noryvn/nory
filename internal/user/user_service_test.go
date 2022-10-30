package user_test

import (
	"context"
	"fmt"
	"testing"

	"nory/domain"
	"nory/internal/class"
	classmember "nory/internal/class_member"
	. "nory/internal/user"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	t.Parallel()
	userRepository := NewUserRepositoryMem()
	classRepository := class.NewClassRepositoryMem()
	classMemberRepository := classmember.NewClassMemberRepositoryMem()

	us := UserService{
		UserRepository:        userRepository,
		ClassRepository:       classRepository,
		ClassMemberRepository: classMemberRepository,
	}

	t.Run("GetUserProfile", func(t *testing.T) {
		t.Parallel()
		user := &domain.User{
			UserId:   uuid.NewString(),
			Name:     "Abelia",
			Username: xid.New().String(),
			Email:    xid.New().String(),
		}
		err := us.UserRepository.CreateUser(context.Background(), user)
		assert.Nil(t, err)

		for i := 0; i < 5; i++ {
			c := &domain.Class{
				OwnerId: user.UserId,
			}
			err := classRepository.CreateClass(context.Background(), c)
			assert.Nil(t, err)

			err = classMemberRepository.CreateMember(context.Background(), &domain.ClassMember{
				UserId:  user.UserId,
				ClassId: c.ClassId,
			})
		}

		res, err := us.GetUserProfile(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, user.Name, res.Data.Name)
		assert.Equal(t, user.UserId, res.Data.UserId)
		assert.Equal(t, 5, res.Data.UserStatistics.JoinedClass)
		assert.Equal(t, 5, res.Data.UserStatistics.OwnedClass)

		_, err = us.UpdateUser(context.Background(), &domain.User{
			UserId:   user.UserId,
			Name:     "Abelia Narindi Agsya",
			Username: "abelia",
		})
		assert.Nil(t, err)

		res, err = us.GetUserProfile(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "Abelia Narindi Agsya", res.Data.Name)
		assert.Equal(t, "abelia", res.Data.Username)
		assert.Equal(t, user.UserId, res.Data.UserId)
		assert.Equal(t, 5, res.Data.UserStatistics.JoinedClass)
		assert.Equal(t, 5, res.Data.UserStatistics.OwnedClass)
	})

	t.Run("GetUserProfileById", func(t *testing.T) {
		t.Parallel()
		user := &domain.User{
			UserId:   uuid.NewString(),
			Name:     "Abelia",
			Username: xid.New().String(),
			Email:    xid.New().String(),
		}
		err := us.UserRepository.CreateUser(context.Background(), user)
		assert.Nil(t, err)

		res, err := us.GetUserProfileById(context.Background(), user.UserId)
		assert.Nil(t, err)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, user.Name, res.Data.Name)
		assert.Equal(t, user.UserId, res.Data.UserId)

		userId := uuid.NewString()
		res, err = us.GetUserProfileById(context.Background(), userId)
		assert.NotNil(t, err)
		msg := fmt.Sprintf("can not find user with id %q", userId)
		assert.Equal(t, msg, err.Error())
	})

	t.Run("GetUserClasses", func(t *testing.T) {
		t.Parallel()
		user := &domain.User{
			UserId:   uuid.NewString(),
			Name:     "Abelia",
			Username: xid.New().String(),
			Email:    xid.New().String(),
		}
		err := us.UserRepository.CreateUser(context.Background(), user)
		assert.Nil(t, err)

		var classes []*domain.Class
		for i := 0; i < 5; i++ {
			class := &domain.Class{
				OwnerId: user.UserId,
				Name:    user.Name,
			}
			classes = append(classes, class)
			err := us.ClassRepository.CreateClass(context.Background(), class)
			assert.Nil(t, err)
		}

		res, err := us.GetUserClasses(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, len(classes), len(res.Data))
		for _, class := range res.Data {
			assert.Equal(t, user.Name, class.Name)
			assert.Equal(t, user.UserId, class.OwnerId)
		}

		upRes, err := us.GetUserProfile(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, 200, upRes.Code)
		assert.Equal(t, user.Name, upRes.Data.Name)
		assert.Equal(t, user.UserId, upRes.Data.UserId)
		assert.Equal(t, len(classes), upRes.Data.UserStatistics.OwnedClass)
	})
}

// testCases := []struct{
// 	user domain.User
// }{
// 	{user: domain.User{Username: strings.Repeat("a", 21)}},
// 	{user: domain.User{Name: strings.Repeat("a", 33)}},
// 	{user: domain.User{Email: strings.Repeat("a", 254) + "@foo.bar"}},
// 	{user: domain.User{Email: strings.Repeat("a", 25)}},
// }

// for _, tc := range testCases {
// 	_, err := us.GetUserProfile(context.Background(), &tc.user)
// 	t.Log(tc.user)
// 	assert.NotNil(t, err)
// }

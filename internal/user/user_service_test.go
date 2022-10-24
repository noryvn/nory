package user_test

import (
	"context"
	"testing"

	"nory/domain"
	"nory/internal/class"
	. "nory/internal/user"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	userRepository := NewUserRepositoryMem()
	classRepository := class.NewClassRepositoryMem()

	us := UserService{
		UserRepository:  userRepository,
		ClassRepository: classRepository,
	}

	t.Run("GetUserProfile", func(t *testing.T) {
		t.Parallel()
		user := &domain.User{
			UserId: uuid.NewString(),
			Name: "Abelia",
			Username: xid.New().String(),
			Email: xid.New().String(),
		}
		err := us.UserRepository.CreateUser(context.Background(), user)
		assert.Nil(t, err)

		res, err := us.GetUserProfile(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, user.Name, res.Data.User.Name)
		assert.Equal(t, user.UserId, res.Data.User.UserId)
	})

	t.Run("GetUserClasses", func(t *testing.T) {
		t.Parallel()
		user := &domain.User{
			UserId: uuid.NewString(),
			Name: "Abelia",
			Username: xid.New().String(),
			Email: xid.New().String(),
		}
		err := us.UserRepository.CreateUser(context.Background(), user)
		assert.Nil(t, err)

		var classes []*domain.Class
		for i := 0; i < 5; i++ {
			class := &domain.Class{
				OwnerId:     user.UserId,
				Name:        "foo",
			}
			classes = append(classes, class)
			err := us.ClassRepository.CreateClass(context.Background(), class)
			assert.Nil(t, err)
		}

		res, err := us.GetUserClasses(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, len(classes), len(res.Data))
		for _, class := range res.Data {
			assert.Equal(t, "foo", class.Name)
			assert.Equal(t, user.UserId, class.OwnerId)
		}

		upRes, err := us.GetUserProfile(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, 200, upRes.Code)
		assert.Equal(t, user.Name, upRes.Data.User.Name)
		assert.Equal(t, user.UserId, upRes.Data.User.UserId)
		assert.Equal(t, len(classes), upRes.Data.OwnedClass)
	})
}

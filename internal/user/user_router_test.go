package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"nory/common/auth"
	"nory/common/response"
	"nory/domain"
	"nory/internal/class"
	classmember "nory/internal/class_member"
	. "nory/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestUserRouter(t *testing.T) {
	t.Parallel()

	userRepository := NewUserRepositoryMem()
	classRepository := class.NewClassRepositoryMem()
	classMemberRepository := classmember.NewClassMemberRepositoryMem()
	classRoute := Route(UserService{
		UserRepository:        userRepository,
		ClassRepository:       classRepository,
		ClassMemberRepository: classMemberRepository,
	})

	app := fiber.New(fiber.Config{
		Immutable: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			e := response.ErrorHandler(c, err)
			if e != nil {
				t.Logf("%#+v", e)
				return e
			}
			return nil
		},
	})
	app.Use(auth.MockMiddleware)
	app.Route("/", classRoute)

	t.Run("unautorization required", func(t *testing.T) {
		for _, tc := range []struct {
			Method string
			Path   string
		}{
			{"GET", "/profile"},
			{"GET", "/classes"},
			{"PATCH", "/user"},
		} {
			req := httptest.NewRequest(tc.Method, tc.Path, nil)
			resp, err := app.Test(req)
			assert.Nil(t, err)
			assert.Equal(t, 401, resp.StatusCode)
		}
	})

	t.Run("profile and classes", func(t *testing.T) {
		user := &domain.User{
			UserId:   uuid.NewString(),
			Name:     xid.New().String(),
			Username: xid.New().String(),
			Email:    xid.New().String(),
		}
		class := &domain.Class{
			OwnerId: user.UserId,
		}

		err := userRepository.CreateUser(context.Background(), user)
		assert.Nil(t, err)
		err = classRepository.CreateClass(context.Background(), class)
		assert.Nil(t, err)
		err = classMemberRepository.CreateMember(context.Background(), &domain.ClassMember{
			UserId:  user.UserId,
			ClassId: class.ClassId,
		})
		assert.Nil(t, err)
		err = classMemberRepository.CreateMember(context.Background(), &domain.ClassMember{
			UserId:  user.UserId,
			ClassId: xid.New().String(),
		})
		assert.Nil(t, err)

		req := httptest.NewRequest("GET", "/profile", nil)
		req.Header.Set("user-id", user.UserId)
		req.Header.Set("username", user.Username)
		req.Header.Set("name", user.Name)
		req.Header.Set("email", user.Email)

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var profile response.Response[*UserProfile]
		err = json.NewDecoder(resp.Body).Decode(&profile)
		assert.Nil(t, err)
		assert.Equal(t, 2, profile.Data.JoinedClass)
		assert.Equal(t, 1, profile.Data.OwnedClass)
		assert.Equal(t, user.Username, profile.Data.User.Username)
		assert.Equal(t, user.Email, profile.Data.User.Email)

		req = httptest.NewRequest("GET", "/classes", nil)
		req.Header.Set("user-id", user.UserId)
		req.Header.Set("username", user.Username)
		req.Header.Set("name", user.Name)
		req.Header.Set("email", user.Email)

		resp, err = app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var classes response.Response[[]*domain.Class]
		err = json.NewDecoder(resp.Body).Decode(&classes)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(classes.Data))

		p := fmt.Sprintf("/id/%s/profile", user.UserId)
		req = httptest.NewRequest("GET", p, nil)

		resp, err = app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var other response.Response[*UserProfile]
		err = json.NewDecoder(resp.Body).Decode(&other)
		assert.Nil(t, err)
		assert.Equal(t, 2, other.Data.JoinedClass)
		assert.Equal(t, 1, other.Data.OwnedClass)
		assert.Equal(t, user.Username, other.Data.User.Username)
		assert.Equal(t, user.Email, other.Data.User.Email)
		assert.Equal(t, profile, other)

		p = fmt.Sprintf("/username/%s/profile", user.Username)
		req = httptest.NewRequest("GET", p, nil)

		resp, err = app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		other = response.Response[*UserProfile]{}
		err = json.NewDecoder(resp.Body).Decode(&other)
		assert.Nil(t, err)
		assert.Equal(t, 2, other.Data.JoinedClass)
		assert.Equal(t, 1, other.Data.OwnedClass)
		assert.Equal(t, user.Username, other.Data.User.Username)
		assert.Equal(t, user.Email, other.Data.User.Email)
		assert.Equal(t, profile, other)

		buff := bytes.NewBuffer(nil)
		err = json.NewEncoder(buff).Encode(domain.User{
			Username: "hai",
		})

		req = httptest.NewRequest("PATCH", "/user", buff)
		req.Header.Set("content-type", "application/json")
		req.Header.Set("user-id", user.UserId)
		req.Header.Set("username", user.Username)
		req.Header.Set("name", user.Name)
		req.Header.Set("email", user.Email)

		resp, err = app.Test(req)
		assert.Nil(t, err)

		p = fmt.Sprintf("/id/%s/profile", user.UserId)
		req = httptest.NewRequest("GET", p, nil)

		resp, err = app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		other = response.Response[*UserProfile]{}
		err = json.NewDecoder(resp.Body).Decode(&other)
		assert.Nil(t, err)
		assert.Equal(t, 2, other.Data.JoinedClass)
		assert.Equal(t, 1, other.Data.OwnedClass)
		assert.Equal(t, "hai", other.Data.User.Username)
		assert.Equal(t, user.Email, other.Data.User.Email)
	})
}

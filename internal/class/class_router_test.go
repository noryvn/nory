package class_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"

	"nory/common/auth"
	"nory/common/response"
	"nory/domain"
	. "nory/internal/class"
	classmember "nory/internal/class_member"
	classtask "nory/internal/class_task"
	"nory/internal/user"
)

func TestClassRouter(t *testing.T) {
	t.Parallel()

	classService := ClassService{
		UserRepository:        user.NewUserRepositoryMem(),
		ClassRepository:       NewClassRepositoryMem(),
		ClassTaskRepository:   classtask.NewClassTaskRepositoryMem(),
		ClassMemberRepository: classmember.NewClassMemberRepositoryMem(),
	}
	classRoute := Route(classService)

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

	t.Run("auth", func(t *testing.T) {
		buff := bytes.NewBuffer(nil)
		err := json.NewEncoder(buff).Encode(&domain.Class{
			Name: "foo-f",
		})
		assert.Nil(t, err)
		userId := uuid.NewString()

		req := httptest.NewRequest("POST", "/create", buff)
		req.Header.Set("content-type", "application/json")
		req.Header.Set("user-id", userId)
		resp, err := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
		var body response.Response[*domain.Class]
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.Nil(t, err)

		class := body.Data

		t.Run("unauthenticated", func(t *testing.T) {
			for _, tc := range []struct {
				Method string
				Path string
			}{
				{"DELETE", fmt.Sprintf("/%s", class.ClassId)},
				{"DELETE", fmt.Sprintf("/%s/member/%s", class.ClassId, userId)},
				{"POST", fmt.Sprintf("/%s/member", class.ClassId)},
				{"POST", fmt.Sprintf("/%s/task", class.ClassId)},
				{"POST", "/create"},
			}{
				req := httptest.NewRequest(tc.Method, tc.Path, nil)
				resp, err := app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 401, resp.StatusCode)
			}
		})
	})

	t.Run("create", func(t *testing.T) {
		for _, tc := range []struct {
			Name string
			User domain.User
			Body domain.Class
			Code int
		}{
			{
				Name: "Success",
				User: domain.User{UserId: uuid.NewString()},
				Body: domain.Class{Name: "foo"},
				Code: 200,
			},
			{
				Name: "unauthenticated",
				User: domain.User{},
				Body: domain.Class{Name: "foo"},
				Code: 401,
			},
		} {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				buff := bytes.NewBuffer(nil)
				err := json.NewEncoder(buff).Encode(tc.Body)
				assert.Nil(t, err)

				req := httptest.NewRequest("POST", "/create", buff)
				req.Header.Set("content-type", "application/json")
				req.Header.Set("user-id", tc.User.UserId)
				resp, err := app.Test(req)
				assert.Equal(t, tc.Code, resp.StatusCode)

				if resp.StatusCode > 299 {
					return
				}
				defer resp.Body.Close()

				var body response.Response[*domain.Class]
				err = json.NewDecoder(resp.Body).Decode(&body)
				assert.Nil(t, err)
				assert.Equal(t, tc.Body.Name, body.Data.Name)
				assert.NotEqual(t, "", body.Data.ClassId)

				p := fmt.Sprintf("/%s/info", body.Data.ClassId)
				req = httptest.NewRequest("GET", p, nil)
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 200, resp.StatusCode)
				err = json.NewDecoder(resp.Body).Decode(&body)
				assert.Nil(t, err)
				assert.Equal(t, tc.Body.Name, body.Data.Name)
				assert.NotEqual(t, "", body.Data.ClassId)

				now := time.Now().UTC()
				p = fmt.Sprintf("/%s/task", body.Data.ClassId)
				buff.Reset()
				err = json.NewEncoder(buff).Encode(domain.ClassTask{
					ClassId:  xid.New().String(),
					AuthorId: uuid.NewString(),
					DueDate:  now.Add(24 * time.Hour),
				})
				assert.Nil(t, err)
				req = httptest.NewRequest("POST", p, buff)
				req.Header.Set("content-type", "application/json")
				req.Header.Set("user-id", tc.User.UserId)
				resp, err = app.Test(req)
				assert.Equal(t, 200, resp.StatusCode)

				req = httptest.NewRequest("GET", p, nil)
				q := req.URL.Query()
				q.Add("from", now.Format(time.RFC3339))
				req.URL.RawQuery = q.Encode()
				resp, err = app.Test(req)
				assert.Equal(t, 200, resp.StatusCode)
				var b response.Response[[]*domain.ClassTask]
				err = json.NewDecoder(resp.Body).Decode(&b)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(b.Data))

				user := &domain.User{
					UserId:   uuid.NewString(),
					Username: xid.New().String(),
					Email:    xid.New().String(),
				}
				err = classService.UserRepository.CreateUser(context.Background(), user)

				buff = bytes.NewBuffer(nil)
				err = json.NewEncoder(buff).Encode(&domain.User{
					Username: user.Username,
				})
				assert.Nil(t, err)
				p = fmt.Sprintf("/%s/member", body.Data.ClassId)
				req = httptest.NewRequest("POST", p, buff)
				req.Header.Set("user-id", tc.User.UserId)
				req.Header.Set("content-type", "application/json")
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 204, resp.StatusCode)

				req = httptest.NewRequest("GET", p, nil)
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 200, resp.StatusCode)

				memBody := response.Response[[]*domain.ClassMember]{}
				err = json.NewDecoder(resp.Body).Decode(&memBody)
				assert.Nil(t, err)
				assert.Equal(t, 2, len(memBody.Data))

				p = fmt.Sprintf("/%s/member/%s", body.Data.ClassId, user.UserId)
				buff.Reset()
				err = json.NewEncoder(buff).Encode(domain.ClassMember{
					Level: "admin",
				})
				req = httptest.NewRequest("PATCH", p, buff)
				req.Header.Set("user-id", tc.User.UserId)
				req.Header.Set("content-type", "application/json")
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 204, resp.StatusCode)

				p = fmt.Sprintf("/%s/member", body.Data.ClassId)
				req = httptest.NewRequest("GET", p, nil)
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 200, resp.StatusCode)

				memBody = response.Response[[]*domain.ClassMember]{}
				err = json.NewDecoder(resp.Body).Decode(&memBody)
				assert.Nil(t, err)
				assert.Equal(t, 2, len(memBody.Data))

				for _, member := range memBody.Data {
					if member.UserId == user.UserId {
						assert.Equal(t, "admin", member.Level)
					}
				}

				p = fmt.Sprintf("/%s/member/%s", body.Data.ClassId, user.UserId)
				req = httptest.NewRequest("DELETE", p, nil)
				req.Header.Set("user-id", tc.User.UserId)
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 204, resp.StatusCode)

				p = fmt.Sprintf("/%s/member", body.Data.ClassId)
				req = httptest.NewRequest("GET", p, nil)
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 200, resp.StatusCode)

				memBody = response.Response[[]*domain.ClassMember]{}
				err = json.NewDecoder(resp.Body).Decode(&memBody)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(memBody.Data))

				p = fmt.Sprintf("/%s", body.Data.ClassId)
				// unauthorized
				req = httptest.NewRequest("DELETE", p, nil)
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 401, resp.StatusCode)
				// unauthenticated
				req = httptest.NewRequest("DELETE", p, nil)
				req.Header.Set("user-id", xid.New().String())
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 403, resp.StatusCode)
				// ok
				req.Header.Set("user-id", tc.User.UserId)
				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 204, resp.StatusCode)

				p = fmt.Sprintf("/%s/info", body.Data.ClassId)
				req = httptest.NewRequest("GET", p, nil)

				resp, err = app.Test(req)
				assert.Nil(t, err)
				assert.Equal(t, 404, resp.StatusCode)
			})
		}
	})
}

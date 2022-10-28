package class_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

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
)

func TestClassRouter(t *testing.T) {
	t.Parallel()

	userRoute := Route(ClassService{
		ClassRepository:       NewClassRepositoryMem(),
		ClassTaskRepository:   classtask.NewClassTaskRepositoryMem(),
		ClassMemberRepository: classmember.NewClassMemberRepositoryMem(),
	})

	app := fiber.New(fiber.Config{
		Immutable: true,
		ErrorHandler: response.ErrorHandler,
	})
	app.Use(auth.MockMiddleware)
	app.Route("/", userRoute)

	t.Run("create", func(t *testing.T) {
		for _, tc := range []struct{
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
		}{
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

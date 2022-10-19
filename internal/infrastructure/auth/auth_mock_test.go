package auth_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"nory/internal/domain"
	. "nory/internal/infrastructure/auth"
	"nory/internal/interfaces"
)

func TestAuthMock(t *testing.T) {
	app := fiber.New(fiber.Config{
		Immutable: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			res, ok := err.(*interfaces.ResponseError)
			if !ok {
				return fiber.DefaultErrorHandler(c, err)
			}
			return c.Status(res.Code).JSON(res)
		},
	})
	app.Use(MockMiddleware)
	app.Get("/", func(c *fiber.Ctx) error {
		u, err := GetUser(c)
		if err != nil {
			return err
		}
		return c.Status(200).JSON(u)
	})

	testCases := []struct {
		Name string
		Code int
		User domain.User
	}{
		{"success", 200, domain.User{UserId: "foo"}},
		{"success", 200, domain.User{UserId: "foo", CreatedAt: time.Now().Round(time.Second)}},
		{"success", 200, domain.User{UserId: "bar", Username: "bar"}},
		{"success", 200, domain.User{UserId: "baz", Username: "baz", Name: "baz"}},
		{"fail", 401, domain.User{UserId: ""}},
		{"fail", 401, domain.User{UserId: "", CreatedAt: time.Now().Round(time.Second)}},
		{"fail", 401, domain.User{UserId: "", Username: "bar"}},
		{"fail", 401, domain.User{UserId: "", Username: "baz", Name: "baz"}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			t.Parallel()

			body := domain.User{}
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("user-id", tc.User.UserId)
			req.Header.Set("username", tc.User.Username)
			req.Header.Set("name", tc.User.Name)
			req.Header.Set("created-at", tc.User.CreatedAt.Format(time.RFC3339))

			resp, err := app.Test(req, 2)
			assert.Equal(t, nil, err, "unexpected error")
			assert.Equal(t, tc.Code, resp.StatusCode, "missmatch status code")

			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.Equal(t, nil, err, "unexpected error, failed to decode body")
			assert.Equal(t, tc.User.UserId, body.UserId, "missmatch user id")
			if tc.Code == 200 {
				assert.Equal(t, tc.User, body, "missmatch user")
			}
		})
	}
}

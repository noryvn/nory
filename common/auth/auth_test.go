package auth_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
	"github.com/stretchr/testify/assert"

	. "nory/common/auth"
	"nory/common/response"
	"nory/internal/user"
)

func TestAuthMiddleware(t *testing.T) {
	supabaseKey := os.Getenv("SUPABASE_KEY")
	supabaseUrl := os.Getenv("SUPABASE_URL")
	email := os.Getenv("AUTH_USER_EMAIL")
	password := os.Getenv("AUTH_USER_PASSWORD")

	if supabaseKey == "" || supabaseUrl == "" || email == "" || password == "" {
		t.Skip()
	}

	userRepository := user.NewUserRepositoryMem()
	supa := supabase.CreateClient(supabaseUrl, supabaseKey)
	a := &Auth{
		SupabaseAuth:   supa.Auth,
		UserRepository: userRepository,
	}

	supaUser, err := a.SupabaseAuth.SignIn(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	assert.Nil(t, err)
	bearer := fmt.Sprintf("Bearer %s", supaUser.AccessToken)

	user, err := a.UserFromBearer(context.Background(), bearer)
	assert.Nil(t, err)
	assert.Equal(t, supaUser.User.ID, user.UserId, "unknown login data received")
	assert.Equal(t, supaUser.User.Email, user.Email, "unknown login data received")

	app := fiber.New(fiber.Config{
		Immutable: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if res, ok := err.(*response.ResponseError); ok {
				return res.Respond(c)
			}
			return fiber.DefaultErrorHandler(c, err)
		},
	})
	app.Use(a.Middleware)
	app.Get("/", func(c *fiber.Ctx) error {
		user, err := GetUser(c)
		if err != nil {
			return err
		}
		return c.JSON(user)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(fiber.HeaderAuthorization, bearer)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	req = httptest.NewRequest("GET", "/", nil)
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

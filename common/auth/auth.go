package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"

	"nory/common/response"
)

var ErrUserNotFound = response.NewUnathorized("authentication required")

const userLocalKey = "authenticated user locals key"

type Auth struct {
	SupabaseAuth *supabase.Auth
}

func (a *Auth) Middleware(c *fiber.Ctx) error {
	bearer := c.Get("authorization")
	if len(bearer) < 7 {
		return response.NewUnathorized("authorization header not valid")
	}
	token := string(bearer[7:])
	user, err := a.SupabaseAuth.User(c.Context(), token)
	if err != nil {
		return err
	}
	c.Locals("user", user)
	return c.Next()
}

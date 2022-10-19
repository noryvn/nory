package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"

	"nory/internal/interfaces"
)

var ErrUserNotFound = interfaces.NewResponseUnathorized("can not found authenticated user")

var userLocalKey = "authenticated user"

type Auth struct {
	SupabaseAuth *supabase.Auth
}

func (a *Auth) Middleware(c *fiber.Ctx) error {
	bearer := c.Get("authorization")
	if len(bearer) < 7 {
		return interfaces.NewResponseUnathorized("authorization header not valid")
	}
	token := string(bearer[7:])
	user, err := a.SupabaseAuth.User(c.Context(), token)
	if err != nil {
		return err
	}
	c.Locals("user", user)
	return c.Next()
}

package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"

	"nory/internal/interfaces"
)

var (
	mock = false
)

func SetMock(m bool) {
	mock = m
}

type Auth struct {
	supabase supabase.Auth
}

func (a *Auth) MiddlewareFiber(c *fiber.Ctx) error {
	bearer := c.Get("authorization")
	if len(bearer) < 7 {
		return interfaces.NewResponseUnathorized("authorization header not valid")
	}
	token := string(bearer[7:])
	user, err := a.supabase.User(c.Context(), token)
	if err != nil {
		return err
	}
	c.Locals("user", user)
	return c.Next()
}

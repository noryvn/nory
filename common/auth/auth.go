package auth

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
	"github.com/rs/xid"

	"nory/common/response"
	"nory/domain"
)

var ErrUserNotFound = response.NewUnathorized("authentication required")

const userLocalKey = "authenticated user locals key"

type Auth struct {
	SupabaseAuth   *supabase.Auth
	UserRepository domain.UserRepository
}

func (a *Auth) Middleware(c *fiber.Ctx) error {
	bearer := c.Get("authorization")
	if len(bearer) < 7 || bearer[:6] != "Bearer" {
		return c.Next()
	}

	user, err := a.UserFromBearer(c.Context(), bearer)
	if err != nil {
		return err
	}

	c.Locals("user", user)

	return c.Next()
}

func (a *Auth) UserFromBearer(ctx context.Context, bearer string) (*domain.User, error) {
	token := string(bearer[7:])
	user, err := a.SupabaseAuth.User(ctx, token)
	if err != nil {
		return nil, err
	}

	u, err := a.UserRepository.GetUser(ctx, user.ID)
	if errors.Is(err, domain.ErrUserNotExists) {
		id := xid.New().String()
		u = &domain.User{
			UserId:   user.ID,
			Username: id,
			Name:     id,
			Email:    user.Email,
		}
		if err := a.UserRepository.CreateUser(ctx, u); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}

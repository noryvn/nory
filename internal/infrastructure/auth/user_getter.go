package auth

import (
	"github.com/gofiber/fiber/v2"

	"nory/internal/domain"
)

func GetUser(c *fiber.Ctx) (*domain.User, error) {
	u, ok := c.Locals(userLocalKey).(*domain.User)
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}

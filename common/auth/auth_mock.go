package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"nory/domain"
)

func MockMiddleware(c *fiber.Ctx) error {
	userId := c.Get("user-id")
	if userId == "" {
		return ErrUserNotFound
	}
	createdAt, _ := time.Parse(time.RFC3339, c.Get("created-at"))
	c.Locals(userLocalKey, &domain.User{
		UserId:    userId,
		CreatedAt: createdAt,
		Username:  c.Get("username"),
		Name:      c.Get("name"),
		Email:     c.Get("email"),
	})
	return c.Next()
}

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
		Username:  c.Get("username"),
		Name:      c.Get("name"),
		CreatedAt: createdAt,
	})
	return c.Next()
}

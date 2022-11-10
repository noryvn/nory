package middleware

import "github.com/gofiber/fiber/v2"

func DefaultHeader(c *fiber.Ctx) error {
	c.Set("cache-control", "private, max-age=0, no-store")
	return c.Next()
}

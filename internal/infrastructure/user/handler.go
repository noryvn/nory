package user

import (
	"nory/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	userRepository domain.UserRepository
}

func NewApp(ur domain.UserRepository) *fiber.App {
	app := fiber.New()
	return app
}

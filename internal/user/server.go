package user

import (
	"nory/common/auth"

	"github.com/gofiber/fiber/v2"
)

type server struct {
	us UserService
}

func CreateApp(userService UserService) *fiber.App {
	s := server{userService}
	app := fiber.New()
	app.Get("/profile", s.GetUserProfile)
	return app
}

func (s server) GetUserProfile(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}
	res, err := s.us.GetUserProfile(c.Context(), user)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

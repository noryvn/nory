package user

import (
	"nory/common/auth"

	"github.com/gofiber/fiber/v2"
)

type server struct {
	us UserService
}

func Route(userService UserService) func(router fiber.Router) {
	s := server{userService}
	return func(router fiber.Router) {
		router.Get("/profile", s.GetUserProfile)
	}
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

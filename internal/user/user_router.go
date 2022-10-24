package user

import (
	"nory/common/auth"

	"github.com/gofiber/fiber/v2"
)

type router struct {
	us UserService
}

func Route(userService UserService) func(router fiber.Router) {
	r := router{userService}
	return func(router fiber.Router) {
		router.Get("/profile", r.GetUserProfile)
		router.Get("/classes", r.GetUserClasses)
	}
}

func (r router) GetUserProfile(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := r.us.GetUserProfile(c.Context(), user)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (r router) GetUserClasses(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := r.us.GetUserClasses(c.Context(), user)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

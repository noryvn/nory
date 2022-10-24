package user

import (
	"nory/common/auth"

	"github.com/gofiber/fiber/v2"
)

type userRouter struct {
	us UserService
}

func Route(userService UserService) func(router fiber.Router) {
	ur := userRouter{userService}
	return func(router fiber.Router) {
		router.Get("/profile", ur.GetUserProfile)
		router.Get("/classes", ur.GetUserClasses)
	}
}

func (ur userRouter) GetUserProfile(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := ur.us.GetUserProfile(c.Context(), user)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (ur userRouter) GetUserClasses(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := ur.us.GetUserClasses(c.Context(), user)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

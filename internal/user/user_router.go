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
		router.Get("/id/:userId/profile", ur.GetOtherUserProfile)
		router.Get("/username/:username/profile", ur.GetOtherUserProfileByUsername)
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

func (ur userRouter) GetOtherUserProfile(c *fiber.Ctx) error {
	userId := c.Params("userId")

	res, err := ur.us.GetUserProfileById(c.Context(), userId)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderCacheControl, "private, max-age=60")
	return res.Respond(c)
}

func (ur userRouter) GetOtherUserProfileByUsername(c *fiber.Ctx) error {
	username := c.Params("username")

	res, err := ur.us.GetUserProfileByUsername(c.Context(), username)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderCacheControl, "private, max-age=60")
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

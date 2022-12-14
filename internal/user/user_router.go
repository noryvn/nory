package user

import (
	"nory/common/auth"
	"nory/domain"

	"github.com/gofiber/fiber/v2"
)

type userRouter struct {
	us UserService
}

func Route(userService UserService) func(router fiber.Router) {
	if userService.UserRepository == nil {
		panic("userRoute: nil UserService.UserRepository")
	}
	if userService.ClassRepository == nil {
		panic("userRoute: nil UserService.ClassRepository")
	}
	if userService.ClassMemberRepository == nil {
		panic("userRoute: nil UserService.ClassMemberRepository")
	}

	ur := userRouter{userService}
	return func(router fiber.Router) {
		router.Get("/profile", ur.GetUserProfile)
		router.Get("/class", ur.GetUserClasses)
		router.Get("/joined", ur.GetUserJoinedClasses)
		router.Get("/id/:userId/profile", ur.GetOtherUserProfile)
		router.Get("/username/:username/profile", ur.GetOtherUserProfileByUsername)
		router.Patch("/profile", ur.PatchUser)
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

func (ur userRouter) PatchUser(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	updateUser := &domain.User{}
	if err := c.BodyParser(updateUser); err != nil {
		return err
	}
	updateUser.UserId = user.UserId

	res, err := ur.us.UpdateUser(c.Context(), updateUser)
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

	return res.Respond(c)
}

func (ur userRouter) GetOtherUserProfileByUsername(c *fiber.Ctx) error {
	username := c.Params("username")

	res, err := ur.us.GetUserProfileByUsername(c.Context(), username)
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

func (ur userRouter) GetUserJoinedClasses(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := ur.us.GetUserJoinedClasses(c.Context(), user)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

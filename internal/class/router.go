package class

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"nory/common/auth"
	"nory/common/response"
	"nory/domain"
)

type userRouter struct {
	cs ClassService
}

func Route(classService ClassService) func(router fiber.Router) {
	ur := userRouter{classService}
	return func(router fiber.Router) {
		router.Get("/:classId/info", ur.classInfo)
		router.Get("/:classId/tasks", ur.classTasks)
		router.Post("/create", ur.create)
	}
}

func (ur userRouter) classInfo(c *fiber.Ctx) error {
	classId := c.Params("classId")
	res, err := ur.cs.GetClassInfo(c.Context(), classId)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

func (ur userRouter) classTasks(c *fiber.Ctx) error {
	var q struct {
		From time.Time
		To   time.Time
	}
	if err := c.QueryParser(&q); err != nil {
		return response.NewBadRequest(err.Error())
	}
	classId := c.Params("classId")
	res, err := ur.cs.GetClassTasks(c.Context(), classId, q.From, q.To)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

func (ur userRouter) create(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}
	class := &domain.Class{
		OwnerId: user.UserId,
	}
	res, err := ur.cs.CreateClass(c.Context(), class)
	if err != nil {
		return err
	}
	return res.Respond(c)
}
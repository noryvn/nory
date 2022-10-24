package class

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"nory/common/auth"
	"nory/common/response"
	"nory/domain"
)

type classRouter struct {
	cs ClassService
}

func Route(classService ClassService) func(router fiber.Router) {
	cr := classRouter{classService}
	return func(router fiber.Router) {
		router.Get("/:classId/info", cr.classInfo)
		router.Get("/:classId/tasks", cr.classTasks)
		router.Post("/create", cr.create)
	}
}

func (cr classRouter) classInfo(c *fiber.Ctx) error {
	classId := c.Params("classId")
	res, err := cr.cs.GetClassInfo(c.Context(), classId)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

func (cr classRouter) classTasks(c *fiber.Ctx) error {
	var q struct {
		From time.Time
		To   time.Time
	}
	if err := c.QueryParser(&q); err != nil {
		return response.NewBadRequest(err.Error())
	}
	classId := c.Params("classId")
	res, err := cr.cs.GetClassTasks(c.Context(), classId, q.From, q.To)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

func (cr classRouter) create(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}
	class := &domain.Class{
		OwnerId: user.UserId,
	}
	res, err := cr.cs.CreateClass(c.Context(), class)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

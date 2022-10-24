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
		router.Get("/:classId/info", cr.getClassInfo)
		router.Get("/:classId/tasks", cr.getClassTasks)
		router.Post("/:classId/task", cr.createClassTask)
		router.Post("/create", cr.createClass)
	}
}

func (cr classRouter) getClassInfo(c *fiber.Ctx) error {
	classId := c.Params("classId")
	res, err := cr.cs.GetClassInfo(c.Context(), classId)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

func (cr classRouter) getClassTasks(c *fiber.Ctx) error {
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

func (cr classRouter) createClassTask(c *fiber.Ctx) error {
	classId := c.Params("classId")

	var task domain.ClassTask
	if err := c.BodyParser(&task); err != nil {
		return err
	}

	task.ClassId = classId
	res, err := cr.cs.CreateClassTask(c.Context(), &task)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) createClass(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	var class domain.Class
	if err := c.BodyParser(&class); err != nil {
		return err
	}

	class.OwnerId = user.UserId
	res, err := cr.cs.CreateClass(c.Context(), &class)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

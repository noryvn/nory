package class

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"nory/common/response"
	"nory/domain"
)

type server struct {
	cs ClassService
}

func Route(classService ClassService) func (router fiber.Router) {
	s := server{classService}
	return func (router fiber.Router) {
		router.Get("/:classId/info", s.ClassInfo)
		router.Get("/:classId/tasks", s.ClassTasks)
		router.Get("/create", s.ClassCreate)
	}
}

func (s server) ClassInfo(c *fiber.Ctx) error {
	classId := c.Params("classId")
	res, err := s.cs.GetClassInfo(c.Context(), classId)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

func (s server) ClassTasks(c *fiber.Ctx) error {
	var q struct {
		From time.Time
		To   time.Time
	}
	if err := c.QueryParser(&q); err != nil {
		return response.NewBadRequest(err.Error())
	}
	classId := c.Params("classId")
	res, err := s.cs.GetClassTasks(c.Context(), classId, q.From, q.To)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

func (s server) ClassCreate(c *fiber.Ctx) error {
	class := &domain.Class{}
	res, err := s.cs.CreateClass(c.Context(), class)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

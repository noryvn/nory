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
	if classService.UserRepository == nil {
		panic("classRoute: nil ClassService.UserRepository")
	}
	if classService.ClassRepository == nil {
		panic("classRoute: nil ClassService.ClassRepository")
	}
	if classService.ClassTaskRepository == nil {
		panic("classRoute: nil ClassService.ClassTaskRepository")
	}
	if classService.ClassMemberRepository == nil {
		panic("classRoute: nil ClassService.ClassMemberRepository")
	}
	if classService.ClassScheduleRepository == nil {
		panic("classRoute: nil ClassService.ClassScheduleRepository")
	}

	cr := classRouter{classService}
	return func(router fiber.Router) {
		router.Delete("/:classId", cr.deleteClass)
		router.Delete("/:classId/member/:memberId", cr.deleteMember)
		router.Delete("/:classId/task/:taskId", cr.deleteClassTask)
		router.Delete("/:classId/schedule/:scheduleId", cr.deleteClassSchedule)
		router.Patch("/:classId/member/:memberId", cr.updateMember)
		router.Patch("/:classId", cr.updateClass)
		router.Get("/:classId/info", cr.getClassInfo)
		router.Get("/info", cr.getClassInfoByName)
		router.Get("/:classId/task", cr.getClassTask)
		router.Get("/:classId/member", cr.listMember)
		router.Get("/:classId/schedule", cr.getClassSchedule)
		router.Post("/:classId/task", cr.createClassTask)
		router.Post("/:classId/schedule", cr.createClassSchedule)
		router.Post("/:classId/member", cr.addMember)
		router.Post("/create", cr.createClass)
	}
}

func (cr classRouter) deleteClass(c *fiber.Ctx) error {
	classId := c.Params("classId")

	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := cr.cs.DeleteClass(c.Context(), user.UserId, classId)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) updateClass(c *fiber.Ctx) error {
	class := &domain.Class{}
	if err := c.BodyParser(class); err != nil {
		return err
	}

	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := cr.cs.UpdateClass(c.Context(), user.UserId, class)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) getClassInfoByName(c *fiber.Ctx) error {
	name := c.Query("name")
	ownerUsername := c.Query("ownerUsername")
	res, err := cr.cs.GetClassInfoByName(c.Context(), ownerUsername, name)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) getClassInfo(c *fiber.Ctx) error {
	classId := c.Params("classId")
	res, err := cr.cs.GetClassInfo(c.Context(), classId)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) getClassTask(c *fiber.Ctx) error {
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

	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	var task domain.ClassTask
	if err := c.BodyParser(&task); err != nil {
		return err
	}

	task.ClassId = classId
	task.AuthorId = user.UserId
	res, err := cr.cs.CreateClassTask(c.Context(), user.UserId, &task)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) deleteClassTask(c *fiber.Ctx) error {
	taskId := c.Params("taskId")

	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := cr.cs.DeleteClassTask(c.Context(), user.UserId, taskId)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) createClassSchedule(c *fiber.Ctx) error {
	classId := c.Params("classId")

	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	var task domain.ClassSchedule
	if err := c.BodyParser(&task); err != nil {
		return err
	}

	task.ClassId = classId
	task.AuthorId = user.UserId
	res, err := cr.cs.CreateSchedule(c.Context(), &task)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) getClassSchedule(c *fiber.Ctx) error {
	classId := c.Params("classId")
	res, err := cr.cs.GetClassSchedules(c.Context(), classId)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) deleteClassSchedule(c *fiber.Ctx) error {
	scheduleId := c.Params("scheduleId")

	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	res, err := cr.cs.DeleteSchedule(c.Context(), user.UserId, scheduleId)
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

func (cr classRouter) addMember(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}
	classId := c.Params("classId")

	var body struct {
		Username string `json:"username"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	res, err := cr.cs.AddMemberByUsername(c.Context(), user.UserId, body.Username, &domain.ClassMember{
		ClassId: classId,
		Level:   "member",
	})
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) listMember(c *fiber.Ctx) error {
	classId := c.Params("classId")

	res, err := cr.cs.ListMember(c.Context(), classId)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) updateMember(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}
	classId := c.Params("classId")
	memberId := c.Params("memberId")

	member := &domain.ClassMember{}
	if err := c.BodyParser(member); err != nil {
		return err
	}
	member.ClassId = classId
	member.UserId = memberId

	res, err := cr.cs.UpdateMember(c.Context(), user.UserId, member)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

func (cr classRouter) deleteMember(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}
	classId := c.Params("classId")
	memberId := c.Params("memberId")

	res, err := cr.cs.DeleteMember(c.Context(), user.UserId, classId, memberId)
	if err != nil {
		return err
	}

	return res.Respond(c)
}

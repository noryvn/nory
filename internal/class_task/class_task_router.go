package classtask

import (
	"nory/common/auth"

	"github.com/gofiber/fiber/v2"
)

type classTaskRouter struct {
	cts ClassTaskService
}

func Route(classTaskService ClassTaskService) func(router fiber.Router) {
	if classTaskService.ClassRepository == nil {
		panic("classTaskRoute: nil classTaskService.ClassRepository")
	}
	if classTaskService.ClassTaskRepository == nil {
		panic("classTaskRoute: nil classTaskService.ClassTaskRepository")
	}

	ctr := classTaskRouter{classTaskService}
	return func(router fiber.Router) {
		router.Delete("/:taskId", ctr.DeleteTask)
	}
}

func (ctr classTaskRouter) DeleteTask(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return err
	}

	taskId := c.Params("taskId")
	res, err := ctr.cts.DeleteTask(c.Context(), taskId, user.UserId)
	if err != nil {
		return err
	}
	return res.Respond(c)
}

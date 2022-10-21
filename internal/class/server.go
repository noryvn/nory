package class

import "github.com/gofiber/fiber/v2"

type server struct {
	cs ClassService
}

func CreateApp(classService ClassService) *fiber.App {
	s := &server{
		cs: classService,
	}
	app := fiber.New()
	app.Get("/info", s.ClassInfo)
	app.Get("/tasks", s.ClassTasks)
	return app
}

func (s *server) ClassInfo(c *fiber.Ctx) error {
	return nil
}

func (s *server) ClassTasks(c *fiber.Ctx) error {
	return nil
}

package class

import "github.com/gofiber/fiber/v2"

type server struct {
	cs ClassService
}

func CreateApp(classService ClassService) *fiber.App {
	app := fiber.New()
	return app
}

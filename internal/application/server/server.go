package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"nory/internal/infrastructure/auth"
	"nory/internal/interfaces"
)

func CreateDevApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: errorHandler,
	})
	useMiddleware(app)
	app.Use(auth.MockMiddleware)
	return app
}

func CreateApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})
	useMiddleware(app)
	return app
}

func useMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(limiter.New(limiter.Config{
		Max: 20,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return interfaces.
				NewResponseTooManyRequest("limit reached, try again in a few minute").
				Respond(c)
		},
	}))
}

func errorHandler(c *fiber.Ctx, err error) error {
	res, ok := err.(*interfaces.ResponseError)
	if !ok {
		return fiber.DefaultErrorHandler(c, err)
	}
	return res.Respond(c)
}

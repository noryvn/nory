package healthcheck

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthCheck struct {
	Pool *pgxpool.Pool
}

func (hc *HealthCheck) Route(r fiber.Router) {
	r.Get("/", hc.Handler)
	r.Head("/", hc.Handler)
}

func (hc *HealthCheck) Handler(c *fiber.Ctx) error {
	if err := hc.Pool.Ping(c.Context()); err != nil {
		fiber.NewError(fiber.StatusInternalServerError, "failed to ping database")
	}

	return nil
}


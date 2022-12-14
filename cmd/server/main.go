package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"

	"nory/common/auth"
	"nory/common/healthcheck"
	"nory/common/middleware"
	"nory/common/response"
	"nory/internal/class"
	"nory/internal/class_member"
	classschedule "nory/internal/class_schedule"
	"nory/internal/class_task"
	"nory/internal/user"
)

func main() {
	addr := getEnv("SERVER_ADDRESS", ":8080")
	allowOrigins := getEnv("ALLOW_ORIGINS", "*")
	dev := getEnv("ENVIRONMENT", "development") == "development"
	databaseUrl := mustGetEnv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	supa := supabase.CreateClient(
		mustGetEnv("SUPABASE_URL"),
		mustGetEnv("SUPABASE_KEY"),
	)

	health := healthcheck.HealthCheck{
		Pool: pool,
	}

	userRepository := user.NewUserRepositoryPostgres(pool)
	classRepository := class.NewClassRepositoryPostgres(pool)
	classTaskRepository := classtask.NewClassTaskRepositoryPostgres(pool)
	classMemberRepository := classmember.NewClassMemberRepositoryPostgres(pool)
	classScheduleRepository := classschedule.NewClassScheduleRepositoryPg(pool)

	userRoute := user.Route(user.UserService{
		UserRepository:        userRepository,
		ClassRepository:       classRepository,
		ClassMemberRepository: classMemberRepository,
	})
	classRoute := class.Route(class.ClassService{
		UserRepository:          userRepository,
		ClassRepository:         classRepository,
		ClassTaskRepository:     classTaskRepository,
		ClassMemberRepository:   classMemberRepository,
		ClassScheduleRepository: classScheduleRepository,
	})
	authMiddleware := auth.Auth{
		SupabaseAuth:   supa.Auth,
		UserRepository: userRepository,
	}

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: dev,
		Immutable:         dev,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			fiberErr, ok := err.(*fiber.Error)
			if ok {
				err = response.NewError(fiberErr.Code, fiberErr.Message)
			}

			err = response.ErrorHandler(c, err)
			if err == nil {
				return nil
			}

			return err
		},
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowHeaders: "*",
		MaxAge:       86400,
	}))
	app.Use(logger.New())
	app.Use(authMiddleware.Middleware)
	app.Use(middleware.DefaultHeader)
	app.Route("/user", userRoute, "user")
	app.Route("/class", classRoute, "class")
	app.Route("/health", health.Route, "health")

	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}

func getEnv(name, def string) string {
	v, exists := os.LookupEnv(name)
	if !exists {
		return def
	}
	return v
}

func mustGetEnv(name string) string {
	v, exists := os.LookupEnv(name)
	if !exists {
		msg := fmt.Sprintf("missing %q in environment variable", name)
		panic(msg)
	}
	return v

}

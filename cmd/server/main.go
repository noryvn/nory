package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedpals/supabase-go"

	"nory/common/auth"
	"nory/common/response"
	"nory/internal/class"
	classtask "nory/internal/class_task"
	"nory/internal/user"
)

func main() {
	addr := getEnv("SERVER_ADDRESS", ":8080")
	dev := getEnv("ENVIRONMENT", "development") == "development"
	databaseUrl := mustGetEnv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		panic(err)
	}

	supa := supabase.CreateClient(
		mustGetEnv("SUPABASE_URL"),
		mustGetEnv("SUPABASE_KEY"),
	)

	userRepository := user.NewUserRepositoryPostgres(pool)
	classRepository := class.NewClassRepositoryPostgres(pool)
	classTaskRepository := classtask.NewClassTaskRepositoryPostgres(pool)

	userRoute := user.Route(user.UserService{
		UserRepository:  userRepository,
		ClassRepository: classRepository,
	})
	classRoute := class.Route(class.ClassService{
		ClassRepository:     classRepository,
		ClassTaskRepository: classTaskRepository,
	})
	authMiddleware := auth.Auth{
		SupabaseAuth:   supa.Auth,
		UserRepository: userRepository,
	}

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: dev,
		Immutable: dev,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			err = response.ErrorHandler(c, err)
			if err == nil {
				return nil
			}
			fmt.Println(err)
			return err
		},
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: !dev,
	}))
	app.Use(authMiddleware.Middleware)
	app.Route("/user", userRoute, "user")
	app.Route("/class", classRoute, "class")

	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}

func getEnv(name, def string) string {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	return v
}

func mustGetEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		msg := fmt.Sprintf("missing %q in environment variable", name)
		panic(msg)
	}
	return v

}

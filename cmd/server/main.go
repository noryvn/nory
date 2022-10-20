package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"nory/internal/user"
)

func main() {
	addr := getEnv("SERVER_ADDRESS", ":8080")

	db, err := pgxpool.New(context.Background(), mustGetEnv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	userRepository := user.NewUserRepositoryPostgres(db)

	userApp := user.CreateApp(user.UserService{
		UserRepository:  userRepository,
		ClassRepository: nil,
	})
	app := fiber.New()
	app.Mount("/user", userApp)
	err = app.Listen(addr)
	if err != nil {
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

package main

import (
	"nory/internal/user"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	addr := getEnv("SERVER_ADDRESS", ":8080")

	userApp := user.CreateApp(user.UserService{
		UserRepository:  nil,
		ClassRepository: nil,
	})
	app := fiber.New()
	app.Mount("/user", userApp)
	err := app.Listen(addr)
	if err != nil {
		panic(err)
	}
}

func getEnv(name, def string) string {
	v := os.Getenv(name)
	if v != "" {
		return def
	}
	return v
}

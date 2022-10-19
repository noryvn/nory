package main

import (
	"nory/internal/application/server"
	"nory/internal/infrastructure/class"
	"nory/internal/infrastructure/user"
)

func main() {
	ur := user.NewUserRepositoryMem()
	cr := class.NewClassRepositoryMem()
	userApp := user.CreateApp(user.UserService{
		UserRepository: ur,
		ClassRepository: cr,
	})
	app := server.CreateDevApp()
	app.Mount("/", userApp)
	app.Listen(":8080")
}

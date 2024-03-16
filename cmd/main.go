package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yaza-putu/crud-fiber/internal/app/todo/handler"
)

func main() {

	app := fiber.New(fiber.Config{
		AppName: "TODO",
	})

	app.Post("/todos", handler.Create)
	app.Put("/todos/:id", handler.Update)
	app.Delete("/todos/:id", handler.Delete)
	app.Get("/todos", handler.All)
	app.Get("/todos/reports", handler.GenReport)
	app.Get("/todos/:id", handler.FindById)

	app.Listen(":3000")
}

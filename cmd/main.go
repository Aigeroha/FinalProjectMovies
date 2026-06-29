package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Fiber v3 успешно запущен!")
	})

	log.Fatal(app.Listen(":8080"))
}

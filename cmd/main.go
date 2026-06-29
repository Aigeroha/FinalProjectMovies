package main

import (
	"log"
	"github.com/gofiber/fiber/v3"
	"final-project/internal/database"
)

func main() {
	
	database.Connect()

	app := fiber.New()

	log.Fatal(app.Listen(":8080"))
}

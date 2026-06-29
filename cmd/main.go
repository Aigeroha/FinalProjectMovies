package main

import (
	"log"
	"github.com/gofiber/fiber/v3"
	"final-project/internal/database"
	"final-project/internal/routes"
)

func main() {
	
	database.Connect()

	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}

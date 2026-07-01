package main

import (
	"log"
	"github.com/gofiber/fiber/v3"
	"final-project/internal/database"
	"final-project/internal/routes"
	"final-project/internal/config"
)

func main() {

	config.LoadConfig()
	
	database.Connect(config.AppConfig.DBConnStr)
	database.ConnectRedis(config.AppConfig.RedisAddr)

	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(config.AppConfig.ServerPort))
}

package main

import (
	"node-backend/database"
	"node-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Create a new Fiber instance
	app := fiber.New()

	// Connect to the database
	database.Connect()

	// Initialize cors
	app.Use(cors.New())

	app.Use(logger.New())

	// Handle routing
	app.Route("/", routes.Router)

	// Listen on port 3000
	app.Listen(":3000")
}

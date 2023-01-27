package main

import (
	"log"
	"node-backend/database"
	"node-backend/routes"
	"node-backend/util/auth"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Connect to the database
	database.Connect()

	app.Use(cors.New())
	app.Use(logger.New())

	// Handle routing
	app.Route("/", routes.Router)

	log.Println(auth.GenerateToken(300))

	// Listen on port 3000
	app.Listen(":3000")
}

package main

import (
	"log"
	"node-backend/database"
	"node-backend/routes"
	"node-backend/util"
	"node-backend/util/auth"
	"time"

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

	token, _ := util.Token(auth.GenerateToken(300), time.Now().Add(time.Hour*24), map[string]interface{}{
		"acc": 123,
		"lvl": 100,
	})

	log.Println("Test token: " + token)

	// Listen on port 3000
	app.Listen(":3000")
}

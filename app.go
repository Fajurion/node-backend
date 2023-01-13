package main

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	// Create a new Fiber instance
	app := fiber.New()

	// Connect to the database
	database.Connect()

	// Initialize cors
	app.Use(cors.New())

	// Handle routing
	app.Route("/", routes.Router)

	// Listen on port 3000
	app.Listen(":3000")
}

func getAccount(c *fiber.Ctx) error {

	id := c.Params("id")

	acc := account.Account{}

	account := database.DBConn.First(&acc, id)

	if account.Error != nil {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(fiber.Map{
		"username": acc.Username,
		"password": acc.Password,
	})
}

func createAccount(ctx *fiber.Ctx) error {
	account := account.Account{
		Username: "test",
		Password: "test123",
	}

	database.DBConn.Create(&account)

	return ctx.SendString("Account created successfully!")
}

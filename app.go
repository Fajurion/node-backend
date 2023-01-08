package main

import (
	"node-backend/database"
	"node-backend/entities/account"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	// Create a new Fiber instance
	app := fiber.New()

	// Connect to the database
	database.Connect()

	// Create a route
	app.Get("/create", createAccount)

	app.Use(cors.New())

	// Listen on port 3000
	app.Listen(":3000")
}

func createAccount(ctx *fiber.Ctx) error {
	account := account.Account{
		Username: "test",
		Password: "test123",
	}

	database.DBConn.Create(&account)

	return ctx.SendString("Account created successfully!")
}

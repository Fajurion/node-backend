package auth

import (
	"node-backend/database"
	"node-backend/entities/account"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Post("/login", login)
	router.Post("/register", register_test)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func login(c *fiber.Ctx) error {
	var req LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid",
		})
	}

	// Get user from database
	var account account.Account
	database.DBConn.Where("email = ?", req.Email).First(&account)

	// Check if user exists
	if account.ID == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid.password",
		})
	}

	// Check if password is correct
	if !account.CheckPassword(req.Password) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid.password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

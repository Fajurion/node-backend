package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Tag      string `json:"tag"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func register_test(c *fiber.Ctx) error {

	// Parse body to register request
	var registerRequest RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid",
		})
	}

	err := database.DBConn.Create(&account.Account{
		Username: registerRequest.Username,
		Tag:      registerRequest.Tag,
		Password: auth.HashPassword(registerRequest.Password),
		Email:    registerRequest.Email,
		Rank:     1, // Default rank
	}).Error

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

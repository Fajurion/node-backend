package auth

import (
	"node-backend/database"
	"node-backend/entities/account"

	"github.com/gofiber/fiber/v2"
)

func register_test(c *fiber.Ctx) error {

	err := database.DBConn.Create(&account.Account{
		Username: "test",
		Tag:      "test",
		Password: "test123",
		Email:    "test@gmail.com",
		Rank:     0,
	}).Error

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

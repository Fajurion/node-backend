package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type registerRequest struct {
	Username string `json:"username"`
	Tag      string `json:"tag"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// When Redis is implemented, this will be replaced with a proper register function.
func register_test(c *fiber.Ctx) error {

	// Parse body to register request
	var registerRequest registerRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return requests.InvalidRequest(c)
	}

	err := database.DBConn.Create(&account.Account{
		ID:       auth.GenerateToken(8),
		Username: registerRequest.Username,
		Tag:      registerRequest.Tag,
		Password: auth.HashPassword(registerRequest.Password),
		Email:    registerRequest.Email,
		RankID:   1, // Default rank
	}).Error

	if err != nil {
		return requests.InvalidRequest(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

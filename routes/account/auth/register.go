package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Tag      string `json:"tag"`
	Password string `json:"password"`
}

// When Redis is implemented, this will be replaced with a proper register function.
func register_test(c *fiber.Ctx) error {

	// Parse body to register request
	var registerRequest registerRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if email is already registered
	valid, normalizedEmail := account.CheckEmail(registerRequest.Email)
	if !valid {
		return requests.FailedRequest(c, "email.registered", nil)
	}

	if database.DBConn.Where("email = ?", normalizedEmail).First(&account.Account{}).RowsAffected > 0 {
		return requests.FailedRequest(c, "email.registered", nil)
	}

	err := database.DBConn.Create(&account.Account{
		ID:       auth.GenerateToken(8),
		Email:    normalizedEmail,
		Username: registerRequest.Username,
		Tag:      registerRequest.Tag,
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

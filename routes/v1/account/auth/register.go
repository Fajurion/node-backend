package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"

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
	if err := util.BodyParser(c, &registerRequest); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if email is already registered
	valid, normalizedEmail := account.CheckEmail(registerRequest.Email)
	if !valid {
		return util.FailedRequest(c, "email.invalid", nil)
	}

	if database.DBConn.Where("email = ?", normalizedEmail).Take(&account.Account{}).RowsAffected > 0 {
		return util.FailedRequest(c, "email.registered", nil)
	}

	var acc account.Account = account.Account{
		ID:       auth.GenerateToken(8),
		Email:    normalizedEmail,
		Username: registerRequest.Username,
		Tag:      registerRequest.Tag,
		RankID:   1, // Default rank
	}

	err := database.DBConn.Create(&acc).Error

	if err != nil {
		return util.InvalidRequest(c)
	}

	err = database.DBConn.Create(&account.Authentication{
		ID:      auth.GenerateToken(8),
		Account: acc.ID,
		Type:    account.TypePassword,
		Secret:  auth.HashPassword(registerRequest.Password),
	}).Error

	if err != nil {
		return util.InvalidRequest(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

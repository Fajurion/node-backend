package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/standards"
	"node-backend/util"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Tag      string `json:"tag"`
	Password string `json:"password"`
	Invite   string `json:"invite"`
}

// When Redis is implemented, this will be replaced with a proper register function.
func register_test(c *fiber.Ctx) error {

	// Parse body to register request
	var req registerRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check the invite
	var invite account.Invite
	if err := database.DBConn.Where("id = ?", req.Invite).Take(&invite).Error; err != nil {
		return util.FailedRequest(c, "invite.invalid", err)
	}

	// Just for security
	if invite.ID != req.Invite {
		return util.FailedRequest(c, "invite.invalid", nil)
	}

	// Check if email is already registered
	valid, normalizedEmail := standards.CheckEmail(req.Email)
	if !valid {
		return util.FailedRequest(c, "email.invalid", nil)
	}

	if database.DBConn.Where("email = ?", normalizedEmail).Take(&account.Account{}).RowsAffected > 0 {
		return util.FailedRequest(c, "email.registered", nil)
	}

	// Check username and tag
	valid, message := standards.CheckUsernameAndTag(req.Username, req.Tag)
	if !valid {
		return util.FailedRequest(c, message, nil)
	}

	var acc account.Account = account.Account{
		ID:       auth.GenerateToken(8),
		Email:    normalizedEmail,
		Username: req.Username,
		Tag:      req.Tag,
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
		Secret:  auth.HashPassword(req.Password),
	}).Error

	if err != nil {
		return util.InvalidRequest(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

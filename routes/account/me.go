package account

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/me
func me(c *fiber.Ctx) error {

	// Get session
	sessionId := util.GetSession(c)

	var session account.Session
	if database.DBConn.Where(&account.Session{ID: sessionId}).Take(&session).Error != nil {
		return requests.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where(&account.Account{ID: session.Account}).Take(&acc).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Retrun details
	return c.JSON(fiber.Map{
		"success": true,
		"account": acc,
	})
}

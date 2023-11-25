package account

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type getRequest struct {
	ID uint `json:"id"`
}

// Route: /account/get
func getAccount(c *fiber.Ctx) error {

	// Parse request
	var req getRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get account
	var account account.Account
	if database.DBConn.Where("id = ?", req.ID).Find(&account).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"name":    account.Username,
		"tag":     account.Tag,
	})
}

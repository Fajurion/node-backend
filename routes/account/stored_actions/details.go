package stored_actions

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type detailsRequest struct {
	Username string `json:"username"`
	Tag      string `json:"tag"`
}

// Route: /account/stored_actions/details
func getDetails(c *fiber.Ctx) error {

	if !util.IsRemoteId(c) {
		return requests.InvalidRequest(c)
	}

	// Parse request
	var req detailsRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where("username = ? AND tag = ?", req.Username, req.Tag).Take(&acc).Error; err != nil {
		return requests.FailedRequest(c, "not.found", err)
	}

	var key account.PublicKey
	if err := database.DBConn.Where("id = ?", acc.ID).Take(&key).Error; err != nil {
		return requests.FailedRequest(c, "not.found", err)
	}

	// Return account details
	return c.JSON(fiber.Map{
		"success": true,
		"account": acc.ID,
		"key":     key.Key,
	})
}

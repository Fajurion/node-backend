package account

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type getRequest struct {
	ID uint `json:"id"`

	// Node data
	Node  uint   `json:"node"`
	Token string `json:"token"`
}

func getAccount(c *fiber.Ctx) error {

	// Parse request
	var req getRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check node token
	_, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Get account
	var account account.Account
	if database.DBConn.Where("id = ?", req.ID).Find(&account).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"account": account,
	})
}

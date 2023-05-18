package connect

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type disconnectRequest struct {
	Node      uint   `json:"node"`
	NodeToken string `json:"token"`
	Session   string `json:"session"`
}

// Route: /node/disconnect
func Disconnect(c *fiber.Ctx) error {

	// Parse request
	var req disconnectRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check node
	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Disconnect account
	if database.DBConn.Model(&account.Session{}).Where("id = ?", req.Session).Update("node", 0).Error != nil {
		log.Println("Failed to disconnect account", req.Session)
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

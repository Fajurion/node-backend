package connect

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type disconnectRequest struct {
	NodeToken string `json:"node_token"`
	Token     string `json:"token"`
}

func Disconnect(c *fiber.Ctx) error {

	// Parse request
	var req disconnectRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	var node node.Node
	if err := database.DBConn.Take(&node, "token = ?", req.NodeToken).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	var session account.Session
	if err := database.DBConn.Take(&session, "token = ?", req.Token).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	session.Device = "app"
	session.Connected = false

	if err := database.DBConn.Save(&session).Error; err != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

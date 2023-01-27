package manage

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type regenRequest struct {
	Node uint `json:"node"` // Node ID
}

func regenToken(c *fiber.Ctx) error {

	// Parse body to remove request
	var req regenRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check permission
	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	if req.Node == 0 {
		return requests.FailedRequest(c, "invalid", nil)
	}

	var node node.Node
	if err := database.DBConn.Take(&node, req.Node).Error; err != nil {
		return requests.FailedRequest(c, "not.found", err)
	}

	node.Token = auth.GenerateToken(300)

	if err := database.DBConn.Save(&node).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return requests.SuccessfulRequest(c)
}

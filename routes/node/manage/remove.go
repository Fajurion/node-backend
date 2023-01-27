package manage

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	Node uint `json:"node"` // Node ID
}

func removeNode(c *fiber.Ctx) error {

	// Parse body to remove request
	var req removeRequest
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

	// Delete node
	if err := database.DBConn.Delete(node.Node{}, req.Node).Error; err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	return requests.SuccessfulRequest(c)
}

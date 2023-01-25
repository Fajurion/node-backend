package manage

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	Token string `json:"token"`

	Node uint `json:"node"` // Node ID
}

func removeNode(c *fiber.Ctx) error {

	// Parse body to remove request
	var req removeRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	// Check if request is valid
	if req.Token == "" || req.Node == 0 {
		return requests.FailedRequest(c, "invalid", nil)
	}

	// Check session and permission
	var session account.Session
	if !requests.CheckSessionPermission(c, req.Token, util.PermissionAdmin, &session) {
		return requests.FailedRequest(c, "no.permission", nil)
	}

	// Delete node
	if err := database.DBConn.Delete(node.Node{}, req.Node).Error; err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	return requests.SuccessfulRequest(c)
}

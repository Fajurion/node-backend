package manage

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type regenRequest struct {
	Token string `json:"token"`

	Node uint `json:"node"` // Node ID
}

func regenToken(c *fiber.Ctx) error {

	// Parse body to remove request
	var req regenRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if req.Token == "" || req.Node == 0 {
		return requests.FailedRequest(c, "invalid", nil)
	}

	// Check session and permission
	var session account.Session
	if !requests.CheckSessionPermission(c, req.Token, util.PermissionManageNodes, &session) {
		return requests.FailedRequest(c, "no.permission", nil)
	}

	var node node.Node
	if err := database.DBConn.First(&node, req.Node).Error; err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	node.Token = auth.GenerateToken(300)

	if err := database.DBConn.Save(&node).Error; err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	return requests.SuccessfulRequest(c)
}

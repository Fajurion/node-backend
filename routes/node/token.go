package node

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"
	"time"

	"github.com/gofiber/fiber/v2"
)

func generateToken(c *fiber.Ctx) error {

	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	tk := auth.GenerateToken(200)

	// Save
	if err := database.DBConn.Create(&node.NodeCreation{
		Token: tk,
		Date:  time.Now(),
	}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   tk,
	})
}

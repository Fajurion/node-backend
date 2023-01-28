package add

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"
	"time"

	"github.com/gofiber/fiber/v2"
)

func createToken(c *fiber.Ctx) error {

	// Check permission
	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	// Create new token
	token := auth.GenerateToken(300)

	if err := database.DBConn.Create(node.NodeCreation{
		Token: token,
		Date:  time.Now(),
	}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
	})
}

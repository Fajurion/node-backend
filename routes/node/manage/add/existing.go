package add

import (
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

func existingNode(c *fiber.Ctx) error {

	// Check permission
	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	return nil
}

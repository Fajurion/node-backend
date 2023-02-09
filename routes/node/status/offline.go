package status

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type offlineRequest struct {
	Token string `json:"token"`
}

func offline(c *fiber.Ctx) error {

	// Parse request
	var req offlineRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get node
	var requested node.Node
	if err := database.DBConn.Where("token = ?", req.Token).Take(&requested).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	// Update status
	requested.Status = node.StatusStopped
	requested.Load = 0

	if err := database.DBConn.Save(&requested).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	return nil
}

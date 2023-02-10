package node

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type thisRequest struct {
	Token string `json:"token"`
}

func this(c *fiber.Ctx) error {

	// Parse request
	var req thisRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get node
	var requested node.Node
	if err := database.DBConn.Where("token = ?", req.Token).Take(&requested).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"node":    requested.ToEntity(),
	})

}

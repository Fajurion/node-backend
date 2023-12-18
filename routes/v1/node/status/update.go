package status

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type updateRequest struct {
	Token     string `json:"token"`
	NewStatus uint   `json:"newStatus"`
}

func update(c *fiber.Ctx) error {

	// Parse request
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get node
	var requested node.Node
	database.DBConn.Where("token = ?", req.Token).Take(&requested)

	if requested.ID == 0 {
		return requests.InvalidRequest(c)
	}

	// Update status
	requested.Status = req.NewStatus
	database.DBConn.Save(&requested)

	return requests.SuccessfulRequest(c)
}

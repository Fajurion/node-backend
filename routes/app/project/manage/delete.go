package manage

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type deleteRequest struct {
	Token   string `json:"token"`
	Project uint   `json:"project"`
}

func deleteProject(c *fiber.Ctx) error {

	// Parse request
	var req deleteRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get node
	var node node.Node
	if err := database.DBConn.Where("token = ?", req.Token).Take(&node).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	// Delete project
	if err := database.DBConn.Where("id = ? AND app = ?", req.Project, node.AppID).Delete(&node).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return requests.SuccessfulRequest(c)
}

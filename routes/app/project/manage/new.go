package manage

import (
	"node-backend/database"
	"node-backend/entities/app/projects"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type newRequest struct {
	Token   string `json:"token"`
	Data    string `json:"data"`
	Creator uint   `json:"creator"`
}

func newProject(c *fiber.Ctx) error {

	// Parse request
	var req newRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get node
	var node node.Node
	if err := database.DBConn.Where("token = ?", req.Token).Take(&node).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	// Create project
	var project projects.Project = projects.Project{
		Creator: req.Creator,
		App:     node.AppID,
		Data:    req.Data,
	}

	if err := database.DBConn.Create(&project).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"project": project,
	})
}

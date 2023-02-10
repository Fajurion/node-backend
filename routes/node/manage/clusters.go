package manage

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type clusterRequest struct {
	Token string `json:"token"`
}

func clusterList(c *fiber.Ctx) error {

	// Parse request
	var req clusterRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	var ct node.NodeCreation
	if err := database.DBConn.Where("token = ?", req.Token).Take(&ct).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	var clusters []node.Cluster
	if err := database.DBConn.Find(&clusters).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"clusters": clusters,
	})
}
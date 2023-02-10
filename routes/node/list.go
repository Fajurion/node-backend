package node

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type listRequest struct {
	Token string `json:"token"`
}

func listNodes(c *fiber.Ctx) error {

	// Get app
	var req listRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	var requested node.Node
	database.DBConn.Where("token = ?", req.Token).Take(&requested)

	if requested.ID == 0 {
		return requests.InvalidRequest(c)
	}

	// Get started nodes
	var nodes []node.Node
	database.DBConn.Where(&node.Node{
		AppID:  requested.AppID,
		Status: node.StatusStarted,
	}).Find(&nodes)

	var startedNodes []node.NodeEntity
	for _, n := range nodes {
		if n.Status == node.StatusStarted && n.ID != requested.ID {
			startedNodes = append(startedNodes, n.ToEntity())
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"nodes":   startedNodes,
	})
}

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

type nodeEntity struct {
	ID     uint   `json:"id"`
	Token  string `json:"token"`
	App    uint   `json:"app"`
	Domain string `json:"domain"`
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
	database.DBConn.Where("app_id = ?", requested.AppID).Find(&nodes)

	var startedNodes []nodeEntity
	for _, n := range nodes {
		if n.Status == node.StatusStarted && n.ID != requested.ID {
			startedNodes = append(startedNodes, nodeEntity{
				ID:     n.ID,
				Token:  n.Token,
				App:    n.AppID,
				Domain: n.Domain,
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"nodes":   startedNodes,
	})
}

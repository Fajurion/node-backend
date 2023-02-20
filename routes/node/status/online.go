package status

import (
	"log"
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type onlineRequest struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

func online(c *fiber.Ctx) error {

	// Parse request
	var req onlineRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get node
	requested, err := nodes.Node(req.ID, req.Token)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Update status
	nodes.TurnOff(requested, node.StatusStarted)

	// Send adoption
	var foundNodes []node.Node
	var startedNodes []node.NodeEntity
	database.DBConn.Model(&node.Node{}).Where("app_id = ?", requested.AppID).Where("status = ?", node.StatusStarted).Find(&foundNodes)

	for _, n := range foundNodes {
		if n.ID != requested.ID {
			if err := n.SendPing(requested); err != nil {

				log.Println("Found offline node: " + n.Domain + "! Shutting down..")

				nodes.TurnOff(requested, node.StatusStopped)
			} else {
				startedNodes = append(startedNodes, n.ToEntity())
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"nodes":   startedNodes,
	})
}

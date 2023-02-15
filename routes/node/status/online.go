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
	Token string `json:"token"`
}

func online(c *fiber.Ctx) error {

	// Parse request
	var req onlineRequest
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
	nodes.TurnOff(requested, node.StatusStarted)

	// Send adoption
	var foundNodes []node.Node
	database.DBConn.Model(&node.Node{}).Where("app_id = ?", requested.AppID).Where("status = ?", node.StatusStarted).Find(&foundNodes)

	for _, n := range foundNodes {
		if n.ID != requested.ID {
			if err := n.SendAdoption(requested); err != nil {

				log.Println("Found offline node: " + n.Domain + "! Shutting down..")

				nodes.TurnOff(requested, node.StatusStopped)
			}
		}
	}

	return requests.SuccessfulRequest(c)
}

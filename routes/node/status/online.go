package status

import (
	"node-backend/database"
	"node-backend/entities/node"
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

	// Send adoption
	var nodes []node.Node
	database.DBConn.Where(&node.Node{
		AppID:  requested.AppID,
		Status: node.StatusStarted,
	}).Find(&nodes)

	for _, n := range nodes {
		if n.ID != requested.ID {
			n.SendAdoption(requested.Domain, requested.Token)
		}
	}

	// Update status
	requested.Status = node.StatusStarted
	database.DBConn.Save(&requested)

	return requests.SuccessfulRequest(c)

}

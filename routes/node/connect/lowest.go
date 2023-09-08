package connect

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type LowestUsageRequest struct {
	Account string `json:"account"`
	Session string `json:"session"`
	Cluster uint   `json:"cluster"`
	App     uint   `json:"app"`
	Node    uint   `json:"node"`  // Node ID
	Token   string `json:"token"` // Node token
}

func GetLowest(c *fiber.Ctx) error {

	// Parse request
	var req LowestUsageRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check node
	_, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Get lowest load node
	var lowest node.Node
	var search node.Node = node.Node{
		ClusterID: req.Cluster,
		AppID:     req.App,
		Status:    node.StatusStarted,
	}

	if err := database.DBConn.Model(&node.Node{}).Where(&search).Order("load DESC").Take(&lowest).Error; err != nil {
		return requests.FailedRequest(c, "not.setup", nil)
	}

	connectionTk, success, err := lowest.GetConnection(req.Account, req.Session, []string{}, node.SenderNode)
	if err != nil {

		if success {
			return requests.FailedRequest(c, err.Error(), nil)
		}

		// Set the node to error
		nodes.TurnOff(&lowest, node.StatusError)

		return requests.FailedRequest(c, "node.error", err)
	}

	// Save node
	if err := database.DBConn.Save(&lowest).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"domain":  lowest.Domain,
		"id":      lowest.ID,
		"token":   connectionTk,
	})
}

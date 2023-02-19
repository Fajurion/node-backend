package connect

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type connectRequest struct {
	Cluster uint   `json:"cluster"`
	App     uint   `json:"app"`
	Token   string `json:"token"`
}

func Connect(c *fiber.Ctx) error {

	// Parse request
	var req connectRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionUseServices) {
		return requests.FailedRequest(c, "no.permission", nil)
	}

	// Get account
	data := util.GetData(c)
	tk := req.Token

	var acc account.Account
	var current account.Session = account.Session{
		Token: "-1",
	}
	if err := database.DBConn.Preload("Sessions").Take(&acc, data["acc"]).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	// Check for too many connections
	connected := 0
	var connectedNode uint
	for _, session := range acc.Sessions {

		if session.Token == tk {
			current = session
		}

		if session.Connected {
			connectedNode = session.Node
			connected++
		}
	}

	if current.Token == "-1" {
		return requests.FailedRequest(c, "not.found", nil)
	}

	if current.Connected {
		return requests.FailedRequest(c, "already.connected", nil)
	}

	if connected >= 3 && !util.Permission(c, util.PermissionServicesUnlimited) {
		return requests.FailedRequest(c, "too.many.connections", nil)
	}

	// Get lowest load node
	var lowest node.Node
	var search node.Node = node.Node{
		ClusterID: req.Cluster,
		AppID:     req.App,
		Status:    node.StatusStarted,
	}

	// Connect to the same node if possible
	if connected > 0 {
		log.Println("Found existing connection!")
		search.ID = connectedNode
	}

	if err := database.DBConn.Model(&node.Node{}).Where(&search).Order("load DESC").Take(&lowest).Error; err != nil {
		return requests.FailedRequest(c, "not.setup", nil)
	}

	connectionTk, err := lowest.GetConnection(util.GetSession(c), uint(data["acc"].(float64)))
	if err != nil {

		// Set the node to error
		nodes.TurnOff(lowest, node.StatusError)

		return requests.FailedRequest(c, "node.error", err)
	}

	current.Connected = true
	current.Node = lowest.ID
	current.App = req.App
	if err := database.DBConn.Save(&current).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
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

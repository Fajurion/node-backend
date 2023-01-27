package connect

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type connectRequest struct {
	Cluster uint `json:"cluster"`
}

func connect(c *fiber.Ctx) error {

	// Parse request
	var req connectRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get account
	data := util.GetData(c)
	tk := data["tk"].(string)

	var acc account.Account
	var current account.Session = account.Session{
		Token: "-1",
	}
	if err := database.DBConn.Preload("Sessions").Take(&acc, data["acc"]).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	// Check for too many connections
	connected := 0
	for _, session := range acc.Sessions {

		if session.Token == tk {
			current = session
		}

		if session.Connected {
			connected++
		}
	}

	if current.Token == "-1" {
		return requests.FailedRequest(c, "not.found", nil)
	}

	if current.Connected {
		return requests.FailedRequest(c, "already.connected", nil)
	}

	if connected >= 3 {
		return requests.FailedRequest(c, "too.many.connections", nil)
	}

	// Check if session is upgraded
	if current.IsWeb() {
		return requests.FailedRequest(c, "not.upgraded", nil)
	}

	// Get lowest load node
	var lowest node.Node
	if err := database.DBConn.Model(&node.Node{}).Where(&node.Node{ClusterID: req.Cluster}).Order("load DESC").Take(&lowest).Error; err != nil {
		return requests.FailedRequest(c, "too.much.load", nil)
	}

	lowest.GetConnection(tk)

	current.Connected = true
	current.Device = "desktop"
	if err := database.DBConn.Save(&current).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"domain":  lowest.Domain,
	})

}

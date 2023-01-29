package connect

import (
	"fmt"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type connectRequest struct {
	Cluster uint `json:"cluster"`
	App     uint `json:"app"`
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
	tk := util.GetToken(c)

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

	if connected >= 3 && !util.Permission(c, util.PermissionServicesUnlimited) {
		return requests.FailedRequest(c, "too.many.connections", nil)
	}

	// Get lowest load node
	var lowest node.Node
	if err := database.DBConn.Model(&node.Node{}).Where(&node.Node{
		ClusterID: req.Cluster,
		AppID:     req.App,
	}).Order("load DESC").Take(&lowest).Error; err != nil {
		return requests.FailedRequest(c, "not.setup", nil)
	}

	connectionTk, err := lowest.GetConnection(tk, uint(data["acc"].(float64)))
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	current.Connected = true
	current.Device = fmt.Sprintf("app:%d", req.App)
	current.Node = lowest.ID
	if err := database.DBConn.Save(&current).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"domain":  lowest.Domain,
		"id":      lowest.ID,
		"token":   connectionTk,
	})
}

package connect

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/nodes"
	"node-backend/util/requests"
	"time"

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
	currentSessionId := util.GetSession(c)
	tk := req.Token

	var acc account.Account
	if err := database.DBConn.Take(&acc, data["acc"]).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	// Get the most recent session
	var mostRecent account.Session
	if err := database.DBConn.Where("account = ?", acc.ID).Order("last_connection DESC").Take(&mostRecent).Error; err != nil {
		mostRecent = account.Session{
			Token: "-1",
		}
	}

	var currentSession account.Session
	if err := database.DBConn.Where("id = ?", currentSessionId).Take(&currentSession).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	if currentSession.Token != tk {
		return requests.FailedRequest(c, "invalid.token", nil)
	}

	if mostRecent.Token == "-1" {
		return requests.FailedRequest(c, "not.found", nil)
	}

	// Get lowest load node
	var lowest node.Node
	var search node.Node = node.Node{
		ClusterID: req.Cluster,
		AppID:     req.App,
		Status:    node.StatusStarted,
	}

	// Connect to the same node if possible
	if mostRecent.Token != "-1" {
		search.ID = mostRecent.Node
	}

	if err := database.DBConn.Model(&node.Node{}).Where(&search).Order("load DESC").Take(&lowest).Error; err != nil {
		return requests.FailedRequest(c, "not.setup", nil)
	}

	connectionTk, success, err := lowest.GetConnection(acc, uint(data["acc"].(float64)))
	if err != nil {

		if success {
			return requests.FailedRequest(c, err.Error(), nil)
		}

		// Set the node to error
		nodes.TurnOff(&lowest, node.StatusError)

		return requests.FailedRequest(c, "node.error", err)
	}

	currentSession.LastConnection = time.Now()
	currentSession.Node = lowest.ID
	currentSession.App = req.App
	if err := database.DBConn.Save(&currentSession).Error; err != nil {
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

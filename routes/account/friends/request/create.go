package request

import (
	"fmt"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/entities/node"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type addFriendRequest struct {
	NodeToken string `json:"node_token"`
	Session   string `json:"session"`
	Username  string `json:"username"`
	Tag       string `json:"tag"`
}

func createRequest(c *fiber.Ctx) error {

	var req addFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	node, err := nodes.Node(req.NodeToken)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	var session account.Session
	if requests.CheckSession(req.Session, &session) {
		return requests.InvalidRequest(c)
	}

	var friend account.Account
	if err := database.DBConn.Where(&account.Account{Username: req.Username, Tag: req.Tag}).Preload("Sessions").Take(&friend).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if the friend is already a friend
	var friendCheck properties.Friend
	if err := database.DBConn.Where(&properties.Friend{Account: friend.ID, Friend: session.Account}).Take(&friendCheck).Error; err == nil && !friendCheck.Request {
		return requests.FailedRequest(c, "already.friends", nil)
	}

	if err := database.DBConn.Where(&properties.Friend{Account: session.Account, Friend: friend.ID}).Take(&properties.Friend{}).Error; err == nil {
		return requests.FailedRequest(c, "already.requested", nil)
	}

	if friendCheck.Request {

		// Accept friend request
		friendCheck.Request = false
		if err := database.DBConn.Save(&friendCheck).Error; err != nil {
			return requests.FailedRequest(c, "server.error", err)
		}

		database.DBConn.Create(&properties.Friend{
			Account: session.Account,
			Friend:  friend.ID,
			Request: false,
		})

		return ExecuteAction(c, "accept", node.ID, node.AppID, friend)
	}

	// Send friend request
	if err := database.DBConn.Create(&properties.Friend{Account: session.Account, Friend: friend.ID, Request: true}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Send notification to friend
	return ExecuteAction(c, "send", node.ID, node.AppID, friend)
}

func ExecuteAction(c *fiber.Ctx, action string, nodeID uint, app uint, friend account.Account) error {

	var session account.Session
	for _, s := range friend.Sessions {
		if s.Connected && s.Device == fmt.Sprintf("app:%d", app) {
			session = s
			break
		}
	}

	if session.Token == "" {
		return requests.SuccessfulRequest(c)
	}

	var nodeInfo node.Node
	if err := database.DBConn.Where(&node.Node{ID: nodeID}).Take(&nodeInfo).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"action":  action,
		"friend":  friend.ID,
		"node":    nodeInfo,
	})

}

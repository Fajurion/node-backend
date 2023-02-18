package request

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type denyFriendRequest struct {
	NodeToken string `json:"node_token"`
	Session   string `json:"session"`
	Account   uint   `json:"username"`
}

// Route: /account/friends/request/deny
func denyRequest(c *fiber.Ctx) error {

	// Parse request
	var req denyFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	node, err := nodes.Node(req.NodeToken)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	var session account.Session
	if requests.CheckSession(req.Session, &session) {
		return requests.FailedRequest(c, "not.found", nil)
	}

	// Check if the friend request exists
	var friendCheck properties.Friend
	if err := database.DBConn.Where(&properties.Friend{Account: session.Account, Friend: req.Account}).Take(&friendCheck).Error; err == nil && !friendCheck.Request {
		return requests.FailedRequest(c, "already.friends", nil)
	}

	var friend account.Account
	if err := database.DBConn.Where(&account.Account{ID: req.Account}).Take(&friend).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	return ExecuteAction(c, "deny", node.ID, node.AppID, friend)
}

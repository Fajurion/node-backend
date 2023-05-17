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
	Node      uint   `json:"id"`
	NodeToken string `json:"token"`
	Session   string `json:"session"`
	Account   string `json:"account"`
}

// Route: /account/friends/request/deny
func denyRequest(c *fiber.Ctx) error {

	// Parse request
	var req denyFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	var session account.Session
	if !requests.GetSession(req.Session, &session) {
		return requests.InvalidRequest(c)
	}

	// Check if there is a friend request
	var friendRequest properties.Friend
	if database.DBConn.Where(&properties.Friend{Account: req.Account, Friend: session.Account, Request: true}).Take(&friendRequest).Error != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	var friendSession account.Session
	database.DBConn.Where(&account.Session{Account: req.Account}).Not("node = ?", 0).Take(&friendSession) // Doesn't matter if the session is connected or not

	// Delete the friend request
	database.DBConn.Delete(&friendRequest)

	return ExecuteAction(c, "deny", req.Account, friendSession, "", "")
}

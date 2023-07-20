package friends

import (
	"fmt"
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type onlineRequest struct {
	Node    uint   `json:"node"`
	Token   string `json:"token"`
	Account string `json:"account"`
}

// Route: /account/friends/online
func onlineFriends(c *fiber.Ctx) error {

	// Parse request
	var req onlineRequest
	if err := c.BodyParser(&req); err != nil {
		requests.DebugRouteError(c, "invalid body: "+err.Error())
		return requests.InvalidRequest(c)
	}

	requests.DebugRouteError(c, fmt.Sprintf("node: %d, token: %s", req.Node, req.Token))

	// Check node
	_, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		requests.DebugRouteError(c, "error getting node: "+err.Error())
		return requests.InvalidRequest(c)
	}

	// Get online friends
	var friends []string
	if err := database.DBConn.Model(&properties.Friend{}).
		Where("account = ? AND EXISTS ( SELECT node FROM sessions WHERE account = friends.friend AND node > 0 )",
			req.Account).
		Pluck("friend", &friends).Error; err != nil {
		requests.DebugRouteError(c, "error getting friends: "+err.Error())
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"friends": friends,
	})
}

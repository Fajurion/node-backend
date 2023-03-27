package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type onlineRequest struct {
	Node      uint   `json:"node"`
	NodeToken string `json:"token"`
	Account   uint   `json:"account"`
}

// Route: /account/friends/online
func onlineFriends(c *fiber.Ctx) error {

	// Parse request
	var req onlineRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check node
	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Get online friends
	var friends []uint
	if err := database.DBConn.Model(&properties.Friend{}).
		Where("account = ? AND EXISTS ( SELECT node FROM sessions WHERE account = friends.friend AND node > 0 )",
			req.Account).
		Pluck("friend", &friends).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"friends": friends,
	})
}

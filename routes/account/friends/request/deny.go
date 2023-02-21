package request

import (
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type denyFriendRequest struct {
	Node      uint   `json:"node"`
	NodeToken string `json:"node_token"`
	Session   uint   `json:"session"`
	Account   uint   `json:"username"`
}

// Route: /account/friends/request/deny
func denyRequest(c *fiber.Ctx) error {

	// Parse request
	var req denyFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// TODO: Deny friend request
	return nil
}

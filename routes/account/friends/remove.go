package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	Node      uint   `json:"node"`
	NodeToken string `json:"token"`
	Account   string `json:"account"`
	Friend    string `json:"friend"`
}

// Route: /account/friends/remove
func removeFriend(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check node
	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Check friend status
	var friendCheck properties.Friend
	if err := database.DBConn.Where(&properties.Friend{Account: req.Account, Friend: req.Friend}).Take(&friendCheck).Error; err != nil {
		return requests.FailedRequest(c, "not.friends", nil)
	}

	if friendCheck.Request {
		return requests.FailedRequest(c, "not.friends", nil)
	}

	// Remove friend
	if database.DBConn.Delete(&friendCheck).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	// Remove friend from other side
	if database.DBConn.Delete(&properties.Friend{Account: req.Friend, Friend: req.Account}).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

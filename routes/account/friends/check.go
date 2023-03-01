package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type checkRequest struct {
	Node      uint   `json:"id"`
	NodeToken string `json:"token"`
	Account   uint   `json:"account"`
	UserIDs   []uint `json:"users"`
}

// Route: /account/friends/check
func checkFriendships(c *fiber.Ctx) error {

	// Parse request
	var req checkRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Check friendships
	var count int64
	if database.DBConn.Where(&properties.Friend{Account: req.Account}).Where("friend IN ?", req.UserIDs).Count(&count).Error != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	if count != int64(len(req.UserIDs)) {
		return requests.FailedRequest(c, "not.friends", nil)
	}

	// Return response
	return requests.SuccessfulRequest(c)
}

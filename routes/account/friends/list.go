package friends

import (
	"node-backend/database"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type friendListRequest struct {
	Request bool  `json:"request"`
	Since   int64 `json:"since"`
}

type friendEntity struct {
	ID        uint   `json:"id"`
	Username  string `json:"name"`
	Tag       string `json:"tag"`
	Signature string `json:"signature"`
	Key       string `json:"key"`
}

func listFriends(c *fiber.Ctx) error {

	// Parse request
	var req friendListRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get token details
	data := util.GetData(c)
	acc := data["acc"].(string)

	// Get friends
	var friends []friendEntity

	// Get requests
	if err := database.DBConn.Table("friends").
		Select("friends.signature, account.id, account.username, account.tag, key.key").
		Where("updated > ? AND account = ? AND request = ?", req.Since, acc, req.Request).
		Joins("join accounts account on account.id = friend").
		Joins("join public_keys key on key.id = friend").Find(&friends).Error; err != nil {

		return requests.FailedRequest(c, "not.found", nil)
	}

	// Return friends
	return c.JSON(fiber.Map{
		"success": true,
		"friends": friends,
	})
}

package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type friendListRequest struct {
	Request bool  `json:"request"`
	Since   int64 `json:"since"`
}

type friendEntity struct {
	Account uint   `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
}

func listFriends(c *fiber.Ctx) error {

	// Parse request
	var req friendListRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get token details
	data := util.GetData(c)
	acc := uint(data["acc"].(float64))

	// Get friends
	var friends []properties.Friend
	if err := database.DBConn.Where(&properties.Friend{Account: acc, Request: req.Request}).Where("updated > ?", req.Since).Preload("AccountData").Find(&friends).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	// Turn into entities
	var friendsEntities []friendEntity
	for _, friend := range friends {
		friendsEntities = append(friendsEntities, friendEntity{
			Account: friend.AccountData.ID,
			Name:    friend.AccountData.Username,
			Tag:     friend.AccountData.Tag,
		})
	}

	// Return friends
	return c.JSON(fiber.Map{
		"success": true,
		"friends": friendsEntities,
	})
}

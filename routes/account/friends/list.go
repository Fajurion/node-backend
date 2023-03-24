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
	Key     string `json:"key"`
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
	if req.Request {

		// Get requests
		if err := database.DBConn.Where(&properties.Friend{Friend: acc, Request: true}).Where("updated > ?", req.Since).Preload("FriendData").Preload("FriendKey").Find(&friends).Error; err != nil {
			return requests.FailedRequest(c, "not.found", nil)
		}
	} else {

		// Get friends
		if err := database.DBConn.Where("updated > ? AND request = ? AND account = ?", req.Since, false, acc).Preload("AccountData").Preload("AccountKey").Find(&friends).Error; err != nil {
			return requests.FailedRequest(c, "not.found", nil)
		}
	}

	// Turn into entities
	var friendsEntities []friendEntity
	for _, friend := range friends {

		if req.Request {
			friend.AccountData = friend.FriendData
			friend.AccountKey = friend.FriendKey
		}

		friendsEntities = append(friendsEntities, friendEntity{
			Account: friend.AccountData.ID,
			Name:    friend.AccountData.Username,
			Tag:     friend.AccountData.Tag,
			Key:     friend.AccountKey.Key,
		})
	}

	// Return friends
	return c.JSON(fiber.Map{
		"success": true,
		"friends": friendsEntities,
	})
}

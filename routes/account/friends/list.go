package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

func listFriends(c *fiber.Ctx) error {

	// Get friends list
	accId := util.GetAcc(c)
	var friends []properties.Friendship
	if err := database.DBConn.Model(&properties.Friendship{}).Where("account = ?", accId).Find(&friends).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Return friends list
	return c.JSON(fiber.Map{
		"success": true,
		"friends": friends,
	})
}

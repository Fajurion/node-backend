package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	Account uint `json:"account"`
}

func removeFriend(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check friend status
	data := util.GetData(c)
	acc := data["acc"].(uint)

	var friendCheck properties.Friend
	if err := database.DBConn.Where(&properties.Friend{Account: acc, Friend: req.Account}).Take(&friendCheck).Error; err != nil {
		return requests.FailedRequest(c, "not.friends", nil)
	}

	// Remove friend
	if err := database.DBConn.Delete(&friendCheck).Error; err != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return nil

}

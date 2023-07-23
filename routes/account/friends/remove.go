package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	ID string `json:"id"`
}

// Route: /account/friends/remove
func removeFriend(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := c.BodyParser(req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Delete friendship
	accId := util.GetAcc(c)
	if err := database.DBConn.Where("id = ? AND account = ?", req.ID, accId).Delete(&properties.Friendship{}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return requests.SuccessfulRequest(c)
}

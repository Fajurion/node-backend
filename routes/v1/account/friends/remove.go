package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type removeRequest struct {
	ID string `json:"id"`
}

// Route: /account/friends/remove
func removeFriend(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if friendship exists
	accId := util.GetAcc(c)
	var friendship properties.Friendship
	if err := database.DBConn.Where("id = ? AND account = ?", req.ID, accId).Take(&friendship).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return requests.FailedRequest(c, "not.found", nil)
		}

		return requests.FailedRequest(c, "server.error", err)
	}

	// Delete friendship
	if err := database.DBConn.Delete(&friendship).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return requests.SuccessfulRequest(c)
}

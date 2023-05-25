package profile

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type editRequest struct {
	Data string `json:"data"`
}

// Route: /account/profile/edit
func editProfile(c *fiber.Ctx) error {

	var req editRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Grab account id from JWT
	id := util.GetAcc(c)

	if database.DBConn.Model(&properties.Profile{}).Where("id = ?", id).Error != nil {
		return requests.FailedRequest(c, "not.enabled", nil)
	}

	if err := database.DBConn.Model(&properties.Profile{}).Where("id = ?", id).Update("data", req.Data).Error; err != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

package profile

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/profile/disable
func disableProfile(c *fiber.Ctx) error {

	var profile properties.Profile
	if err := database.DBConn.Where("id = ?", util.GetAcc(c)).First(&profile).Error; err != nil {
		return requests.FailedRequest(c, "not.enabled", nil)
	}

	if err := database.DBConn.Delete(&profile).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return requests.SuccessfulRequest(c)
}

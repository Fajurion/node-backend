package profile

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/profile/me
func currentProfile(c *fiber.Ctx) error {

	accId := util.GetAcc(c)

	var profile properties.Profile
	if database.DBConn.Where("id = ?", accId).Take(&profile).Error != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"enabled": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"enabled": true,
		"profile": profile,
	})
}

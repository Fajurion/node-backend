package stored_actions

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

func listStoredActions(c *fiber.Ctx) error {

	// Get stored actions
	accId := util.GetAcc(c)
	var storedActions []properties.StoredAction
	if database.DBConn.Where("account = ?", accId).Find(&storedActions).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	// Return stored actions
	return c.JSON(fiber.Map{
		"success": true,
		"actions": storedActions,
	})
}

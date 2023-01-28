package manage

import (
	"node-backend/database"
	"node-backend/entities/app"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	ID uint `json:"id"`
}

func removeApp(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	// Delete app
	if err := database.DBConn.Delete(&app.App{}, req.ID).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	// TOOD: Purge everything related to the app

	return requests.SuccessfulRequest(c)
}

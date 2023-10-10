package app

import (
	"node-backend/database"
	"node-backend/entities/app"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type getVersionRequest struct {
	App uint `json:"app"`
}

// Route: /app/version
func getVersion(c *fiber.Ctx) error {

	var req getVersionRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	var app app.App
	if database.DBConn.Where("id = ?", req.App).Take(&app).Error != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"version": app.Version,
	})
}

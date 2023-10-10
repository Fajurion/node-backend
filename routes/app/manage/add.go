package manage

import (
	"node-backend/database"
	"node-backend/entities/app"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type addRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AccessLevel uint   `json:"access_level"`
}

// Route: /app/manage/add
func addApp(c *fiber.Ctx) error {

	// Parse request
	var req addRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if len(req.Name) < 3 || len(req.Description) < 3 {
		return requests.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	// Create app
	created := app.App{
		Name:        req.Name,
		Description: req.Description,
		Version:     0,
		AccessLevel: req.AccessLevel,
	}

	if err := database.DBConn.Create(&created).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"app":     created,
	})
}

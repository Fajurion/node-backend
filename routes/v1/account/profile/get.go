package profile

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type getProfileRequest struct {
	ID string `json:"id"`
}

// Route: /account/profile/get
func getProfile(c *fiber.Ctx) error {

	var req getProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	var profile properties.Profile
	if err := database.DBConn.Where("id = ?", req.ID).Take(&profile).Error; err != nil {
		return requests.FailedRequest(c, requests.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"profile": profile,
	})
}

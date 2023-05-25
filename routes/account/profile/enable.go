package profile

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type enableRequest struct {
	Key string `json:"key"`
}

// Route: /account/profile/enable
func enableProfile(c *fiber.Ctx) error {

	// Parse request
	var req enableRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Grab account id from JWT
	id := util.GetAcc(c)

	// Check if profile is already enabled
	if database.DBConn.Model(&properties.Profile{}).Where("id = ?", id).Error == nil {
		return requests.FailedRequest(c, "already.enabled", nil)
	}

	// Create profile
	profile := properties.Profile{
		ID:   id,
		Key:  req.Key,
		Data: "",
	}

	if err := database.DBConn.Create(&profile).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Return success
	return requests.SuccessfulRequest(c)
}

package profile

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type getProfileRequest struct {
	ID string `json:"id"`
}

// Route: /account/profile/get
func getProfile(c *fiber.Ctx) error {

	var req getProfileRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	var profile properties.Profile
	if err := database.DBConn.Where("id = ?", req.ID).Take(&profile).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"profile": profile,
	})
}

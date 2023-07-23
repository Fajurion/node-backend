package profile

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type getRequest struct {
	Account string `json:"account"`
}

// Route: /account/profile/get
func getProfile(c *fiber.Ctx) error {

	var req getRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	//accId := util.GetAcc(c)

	// TODO: Reimplement friendship check
	/*
		if database.DBConn.Model(&properties.Friend{}).Where("account = ? AND friend = ? AND request = ?", accId, req.Account, false).
			Take(&properties.Friend{}).Error != nil {
			return requests.FailedRequest(c, "not.friends", nil)
		}
	*/

	var profile properties.Profile
	if err := database.DBConn.Where("id = ?", req.Account).Take(&profile).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"profile": profile,
	})
}

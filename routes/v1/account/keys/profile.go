package keys

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/profile/get
func getProfileKey(c *fiber.Ctx) error {

	// Get account
	accId := util.GetAcc(c)

	// Get public key
	var key account.ProfileKey
	if database.DBConn.Where("id = ?", accId).Take(&key).Error != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	if key.Key == "" {
		return requests.FailedRequest(c, "not.found", nil)
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"key":     key.Key,
	})
}

// Route: /account/keys/profile/set
func setProfileKey(c *fiber.Ctx) error {

	var req setRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get account
	accId := util.GetAcc(c)

	var acc account.Account
	if database.DBConn.Where("id = ?", accId).Take(&acc).Error != nil {
		return requests.InvalidRequest(c)
	}

	if database.DBConn.Where("id = ?", accId).Take(&account.ProfileKey{}).Error == nil {
		return requests.FailedRequest(c, "already.set", nil)
	}

	// Set public key
	if database.DBConn.Create(&account.ProfileKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

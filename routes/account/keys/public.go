package keys

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/public/get
func getPublicKey(c *fiber.Ctx) error {

	// Get account
	accId := util.GetAcc(c)

	// Get public key
	var key account.PublicKey
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

type setRequest struct {
	Password string `json:"password"`
	Key      string `json:"key"`
}

// Route: /account/keys/public/set
func setPublicKey(c *fiber.Ctx) error {

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

	// Check password
	if !acc.CheckPassword(req.Password) {
		return requests.FailedRequest(c, "invalid.password", nil)
	}

	// Set public key
	database.DBConn.Where("id = ?", accId).Delete(&account.PublicKey{})
	if database.DBConn.Create(&account.PublicKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

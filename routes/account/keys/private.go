package keys

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/private/get
func getPrivateKey(c *fiber.Ctx) error {

	// Get account
	accId := util.GetAcc(c)

	// Get private key
	var key account.PrivateKey
	if database.DBConn.Where("id = ?", accId).Take(&key).Error != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"key":     key.Key,
	})
}

type setPrivateKeyRequest struct {
	Key string `json:"key"`
}

// Route: /account/keys/private/set
func setPrivateKey(c *fiber.Ctx) error {

	// Parse request
	var req setPrivateKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get account
	accId := util.GetAcc(c)

	var acc account.Account
	if database.DBConn.Where("id = ?", accId).Take(&acc).Error != nil {
		return requests.InvalidRequest(c)
	}

	if database.DBConn.Where("id = ?", accId).Take(&account.PrivateKey{}).Error == nil {
		return requests.FailedRequest(c, "already.set", nil)
	}

	// Set private key
	database.DBConn.Where("id = ?", accId).Delete(&account.PrivateKey{})
	if database.DBConn.Create(&account.PrivateKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

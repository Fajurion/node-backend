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
	data := util.GetData(c)
	accId := uint(data["acc"].(float64))

	// Get private key
	var key account.PrivateKey
	if database.DBConn.Find(&key, accId).Error != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"key":     key.Key,
	})
}

type setPrivateKeyRequest struct {
	Password string `json:"password"`
	Key      string `json:"key"`
}

// Route: /account/keys/private/set
func setPrivateKey(c *fiber.Ctx) error {

	// Parse request
	var req setPrivateKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	// Get account
	data := util.GetData(c)
	accId := data["acc"].(string)

	var acc account.Account
	if database.DBConn.Find(&acc, accId).Error != nil {
		return requests.InvalidRequest(c)
	}

	// Check password
	if !acc.CheckPassword(req.Password) {
		return requests.FailedRequest(c, "invalid.password", nil)
	}

	// Set private key
	database.DBConn.Delete(&account.PrivateKey{}, accId)
	if database.DBConn.Create(&account.PrivateKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

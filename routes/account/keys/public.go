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
	data := util.GetData(c)
	accId := uint(data["acc"].(float64))

	// Get public key
	var key account.PublicKey
	if database.DBConn.Find(&key, accId).Error != nil {
		return requests.InvalidRequest(c)
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

	// Set public key
	database.DBConn.Delete(&account.PublicKey{}, accId)
	if database.DBConn.Create(&account.PublicKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}

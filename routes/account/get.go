package account

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type getRequest struct {
	ID string `json:"id"`
}

// Route: /account/get
func getAccount(c *fiber.Ctx) error {

	// Parse request
	var req getRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println(err)
		return requests.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Select("username", "tag").Where("id = ?", req.ID).Take(&acc).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	var pub account.PublicKey
	if err := database.DBConn.Select("key").Where("id = ?", req.ID).Take(&pub).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	var signaturePub account.SignatureKey
	if err := database.DBConn.Select("key").Where("id = ?", req.ID).Take(&signaturePub).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"name":    acc.Username,
		"tag":     acc.Tag,
		"sg":      signaturePub.Key,
		"pub":     pub.Key,
	})
}

package vault

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/vault/list
func listEntries(c *fiber.Ctx) error {

	// Get friends list
	accId := util.GetAcc(c)
	var entries []properties.VaultEntry
	if err := database.DBConn.Model(&properties.VaultEntry{}).Where("account = ?", accId).Find(&entries).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Return friends list
	return c.JSON(fiber.Map{
		"success": true,
		"friends": entries,
	})
}
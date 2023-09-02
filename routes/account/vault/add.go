package vault

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type addEntryRequest struct {
	Tag     string `json:"tag"`     // Tag
	Payload string `json:"payload"` // Encrypted payload
}

// Route: /account/vault/add
func addEntry(c *fiber.Ctx) error {

	// Parse request
	var req addEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if the account has too many entries
	accId := util.GetAcc(c)
	var entryCount int64
	if err := database.DBConn.Model(&properties.VaultEntry{}).Where("account = ?", accId).Count(&entryCount).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	if entryCount >= MaximumEntries {
		return requests.FailedRequest(c, "limit.reached", nil)
	}

	// Create vault entry
	vaultEntry := properties.VaultEntry{
		ID:      auth.GenerateToken(12),
		Account: accId,
		Tag:     req.Tag,
		Payload: req.Payload,
	}
	if err := database.DBConn.Create(&vaultEntry).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      vaultEntry.ID,
	})
}

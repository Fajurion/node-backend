package vault

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type existsRequest struct {
	Hash string `json:"hash"`
}

// Route: /account/vault/exists
func existsEntry(c *fiber.Ctx) error {

	var req existsRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if the friendship exists
	accId := util.GetAcc(c)
	if err := database.DBConn.Where("account = ? AND hash = ?", accId, req.Hash).Take(&properties.VaultEntry{}).Error; err != nil {
		return requests.FailedRequest(c, "not.found", nil)
	}

	return requests.SuccessfulRequest(c)
}

package account

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

// generateRemoteId generates a remote id (Route: /account/remote_id)
func generateRemoteId(c *fiber.Ctx) error {

	acc := util.GetAcc(c)
	var account account.Account
	if database.DBConn.Where("id = ?", acc).Preload("Rank").Take(&account).Error != nil {
		return util.InvalidRequest(c)
	}

	// Generate remote id
	rid, err := util.RemoteId(account.Rank.Level, auth.GenerateToken(32))
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"id":      rid,
	})
}

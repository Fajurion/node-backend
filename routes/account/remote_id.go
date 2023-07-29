package account

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

// generateRemoteId generates a remote id (Route: /account/remote_id)
func generateRemoteId(c *fiber.Ctx) error {

	acc := util.GetAcc(c)
	var account account.Account
	if database.DBConn.Where("id = ?", acc).Preload("Rank").Take(&account).Error != nil {
		return requests.InvalidRequest(c)
	}

	// Generate remote id
	rid, err := util.RemoteId(account.Rank.Level, auth.GenerateToken(32))
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      rid,
	})
}

package session

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

func upgrade(c *fiber.Ctx) error {

	// Check session
	data := util.GetData(c)
	tk := data["tk"].(string)

	var session account.Session
	if requests.CheckSession(c, tk, &session) {
		return requests.InvalidRequest(c)
	}

	if session.IsDesktop() {
		return requests.FailedRequest(c, "already.upgraded", nil)
	}

	// Upgrade session
	session.Upgrade()
	err := database.DBConn.Save(&session).Error

	if err != nil {
		return requests.FailedRequest(c, "failed.upgrade", err)
	}

	return requests.SuccessfulRequest(c)
}

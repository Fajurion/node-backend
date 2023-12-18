package session

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

func logOut(c *fiber.Ctx) error {

	// Get token
	sessionId := util.GetSession(c)

	var session account.Session
	if !util.GetSession(sessionId, &session) {
		return util.InvalidRequest(c)
	}

	// Log out
	database.DBConn.Delete(&session)

	return util.SuccessfulRequest(c)
}

package session

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

func logOut(c *fiber.Ctx) error {

	// Get token
	token := util.GetData(c)

	var session account.Session
	if requests.CheckSession(token["tk"].(string), &session) {
		return requests.InvalidRequest(c)
	}

	// Log out
	database.DBConn.Delete(&session)

	return requests.SuccessfulRequest(c)
}

package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"
	"time"

	"github.com/gofiber/fiber/v2"
)

func refreshSession(c *fiber.Ctx) error {

	// Check if session is valid
	data := util.GetData(c)
	tk := data["tk"].(string)

	var session account.Session
	if requests.CheckSession(c, tk, &session) {
		return requests.InvalidRequest(c)
	}

	if session.Account != data["acc"] {
		return requests.InvalidRequest(c)
	}

	if session.IsExpired() {
		database.DBConn.Delete(&session)
		return requests.FailedRequest(c, "session.expired", nil)
	}

	// Check session duration
	if time.Until(session.End) > time.Hour*16*3 {
		return requests.FailedRequest(c, "session.duration", nil)
	}

	// Refresh session
	session.End = time.Now().Add(time.Hour * 24 * 3)
	database.DBConn.Save(&session)

	jwtToken, err := util.Token(tk, time.Now().Add(time.Hour*24), data)
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   jwtToken,
	})
}

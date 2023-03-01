package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"
	"time"

	"github.com/gofiber/fiber/v2"
)

type refreshRequest struct {
	Token string `json:"token"`
}

// Route: /auth/refresh
func refreshSession(c *fiber.Ctx) error {

	// Parse request
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if session is valid
	data := util.GetData(c)

	var session account.Session
	if requests.CheckSession(req.Token, &session) {
		return requests.InvalidRequest(c)
	}

	if session.Account != data["acc"] {
		return requests.InvalidRequest(c)
	}

	if session.IsExpired() {
		database.DBConn.Delete(&session)
		return requests.FailedRequest(c, "session.expired", nil)
	}

	// Refresh session
	session.End = time.Now().Add(time.Hour * 24 * 7)
	database.DBConn.Save(&session)

	// Check session duration
	if time.Until(session.End) > time.Hour*20*7 {
		return requests.FailedRequest(c, "session.duration", nil)
	}

	jwtToken, err := util.Token(session.ID, time.Now().Add(time.Hour*24*3), data)
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"token":         jwtToken,
		"refresh_token": session.Token,
	})
}

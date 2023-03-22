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
	Session uint   `json:"session"`
	Token   string `json:"token"`
}

// Route: /auth/refresh
func refreshSession(c *fiber.Ctx) error {

	// Parse request
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if session is valid
	var session account.Session
	if !requests.GetSession(req.Session, &session) {
		return requests.InvalidRequest(c)
	}

	if session.Token != req.Token {
		return requests.InvalidRequest(c)
	}

	// Refresh session
	session.LastUsage = time.Now().Add(time.Hour * 24 * 7)
	database.DBConn.Save(&session)

	// Create new token
	jwtToken, err := util.Token(session.ID, time.Now().Add(time.Hour*24*3), fiber.Map{
		"acc": session.Account,
		"lvl": session.PermissionLevel,
	})

	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"token":         jwtToken,
		"refresh_token": session.Token,
	})
}

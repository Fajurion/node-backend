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
	Session string `json:"session"`
	Token   string `json:"token"`
}

// Route: /auth/refresh
func refreshSession(c *fiber.Ctx) error {

	// Parse request
	var req refreshRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if session is valid
	var session account.Session
	if !requests.GetSession(req.Session, &session) {
		return util.InvalidRequest(c)
	}

	if session.Token != req.Token {
		return util.InvalidRequest(c)
	}

	// Refresh session
	session.LastUsage = time.Now().Add(time.Hour * 24 * 7)
	database.DBConn.Save(&session)

	// Create new token
	jwtToken, err := util.Token(session.ID, session.Account, session.PermissionLevel, time.Now().Add(time.Hour*24*3))

	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"token":         jwtToken,
		"refresh_token": session.Token,
	})
}

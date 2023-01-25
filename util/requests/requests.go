package requests

import (
	"node-backend/database"
	"node-backend/entities/account"

	"github.com/gofiber/fiber/v2"
)

func SuccessfulRequest(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"success": true,
	})
}

func FailedRequest(c *fiber.Ctx, error string, err error) error {
	return c.Status(400).JSON(fiber.Map{
		"success": false,
		"error":   error,
	})
}

// CheckSession checks if the session is valid (returns true if it isn't)
func CheckSession(c *fiber.Ctx, token string, session account.Session) bool {

	err := database.DBConn.First(session, token).Error
	return err != nil
}

// CheckSessionPermission checks if the session has the required permission level (returns true if it doesn't)
func CheckSessionPermission(c *fiber.Ctx, token string, permission uint, session *account.Session) bool {

	err := database.DBConn.First(session, token).Error

	if err != nil || session.PermissionLevel < permission {
		return true
	}

	return false
}

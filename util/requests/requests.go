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
	return c.Status(200).JSON(fiber.Map{
		"success": false,
		"error":   error,
	})
}

func InvalidRequest(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusBadRequest)
}

// GetSession gets the session from the database (returns false if it doesn't exist)
func GetSession(id string, session *account.Session) bool {

	if err := database.DBConn.Take(&session, id).Error; err != nil {
		return false
	}

	return true
}

// CheckSessionPermission checks if the session has the required permission level (returns true if it doesn't)
func CheckSessionPermission(token string, permission uint, session *account.Session) bool {

	err := database.DBConn.Take(session, token).Error

	if err != nil || session.PermissionLevel < permission {
		return true
	}

	return false
}

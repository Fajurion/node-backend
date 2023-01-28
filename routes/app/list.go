package app

import (
	"node-backend/database"
	"node-backend/entities/app"

	"github.com/gofiber/fiber/v2"
)

func listApps(c *fiber.Ctx) error {

	var apps []app.App
	database.DBConn.Find(&apps)

	return c.JSON(apps)
}

package util

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

const ErrorNode = "node.error"
const ErrorServer = "server.error"

func DebugRouteError(c *fiber.Ctx, msg string) {
	if Testing {
		log.Println(c.Route().Path+":", msg)
	}
}

func SuccessfulRequest(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"success": true,
	})
}

func FailedRequest(c *fiber.Ctx, error string, err error) error {

	if LogErrors && err != nil {
		log.Println(c.Route().Path+":", err)
	}

	return c.Status(200).JSON(fiber.Map{
		"success": false,
		"error":   error,
	})
}

func InvalidRequest(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusBadRequest)
}

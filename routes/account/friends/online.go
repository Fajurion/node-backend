package friends

import "github.com/gofiber/fiber/v2"

func onlineFriends(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
	})
}

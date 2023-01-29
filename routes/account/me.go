package account

import "github.com/gofiber/fiber/v2"

func me(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

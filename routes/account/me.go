package account

import "github.com/gofiber/fiber/v2"

// Route: /account/me
func me(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

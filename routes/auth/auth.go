package auth

import "github.com/gofiber/fiber/v2"

func Setup(router fiber.Router) {
	router.Post("/login", login)
}

func login(c *fiber.Ctx) error {
	return c.SendString("Login route")
}


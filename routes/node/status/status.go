package status

import "github.com/gofiber/fiber/v2"

func Setup(router fiber.Router) {
	router.Post("/update", update)
}

package request

import "github.com/gofiber/fiber/v2"

func Setup(router fiber.Router) {
	router.Post("/create", createRequest)
	router.Post("/deny", denyRequest)
}

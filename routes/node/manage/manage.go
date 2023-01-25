package manage

import "github.com/gofiber/fiber/v2"

func Setup(router fiber.Router) {
	router.Post("/add", addNode)
	router.Post("/remove", removeNode)
	router.Post("/regen", regenToken)
}

package account

import "github.com/gofiber/fiber/v2"

func Authorized(router fiber.Router) {
	router.Post("/list", listProjects)
}

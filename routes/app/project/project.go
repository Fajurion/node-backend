package project

import "github.com/gofiber/fiber/v2"

func Unauthorized(router fiber.Router) {
	router.Post("/fetch", fetch)
}

func Authorized(router fiber.Router) {
}

package auth

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/login", login)
	router.Post("/register", register_test)
}

func Authorized(router fiber.Router) {
	router.Post("/refresh", refreshSession)
}

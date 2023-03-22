package auth

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/login", login)
	router.Post("/register", register_test)
	router.Post("/refresh", refreshSession)
}

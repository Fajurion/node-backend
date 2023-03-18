package keys

import "github.com/gofiber/fiber/v2"

func Authorized(router fiber.Router) {
	router.Post("/public/get", getPublicKey)
	router.Post("/public/set", setPublicKey)
}

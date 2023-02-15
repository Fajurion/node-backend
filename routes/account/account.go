package account

import (
	"node-backend/routes/account/friends"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/friends", friends.Unauthorized)
}

func Authorized(router fiber.Router) {
	router.Post("/me", me)
	router.Route("/friends", friends.Authorized)
}

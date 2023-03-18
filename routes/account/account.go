package account

import (
	"node-backend/routes/account/friends"
	"node-backend/routes/account/keys"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/friends", friends.Unauthorized)
}

func Authorized(router fiber.Router) {
	router.Post("/me", me)
	router.Route("/friends", friends.Authorized)
	router.Route("/keys", keys.Authorized)
}

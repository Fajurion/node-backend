package account

import (
	"node-backend/routes/account/friends"
	"node-backend/routes/account/keys"
	"node-backend/routes/account/rank"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/friends", friends.Unauthorized)
	router.Route("/rank", rank.Unauthorized)

	router.Post("/get", getAccount)
}

func Authorized(router fiber.Router) {
	router.Route("/friends", friends.Authorized)
	router.Route("/keys", keys.Authorized)

	router.Post("/me", me)
}

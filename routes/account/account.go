package account

import (
	"node-backend/routes/account/keys"
	"node-backend/routes/account/profile"
	"node-backend/routes/account/rank"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/rank", rank.Unauthorized)

	router.Post("/get", getAccount)
}

func Authorized(router fiber.Router) {
	router.Route("/keys", keys.Authorized)
	router.Route("/profile", profile.SetupRoutes)

	router.Post("/remote_id", generateRemoteId)
	router.Post("/me", me)
}

package account

import (
	"node-backend/routes/account/friends"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/friends", friends.Setup)
	router.Post("/me", me)
}

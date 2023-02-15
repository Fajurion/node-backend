package friends

import (
	"node-backend/routes/account/friends/request"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/request", request.Setup)
}

func Authorized(router fiber.Router) {
	router.Post("/remove", removeFriend)
}

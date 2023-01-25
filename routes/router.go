package routes

import (
	"node-backend/routes/account/auth"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	router.Route("/auth", auth.Setup)
}

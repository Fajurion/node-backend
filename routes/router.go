package routes

import (
	"node-backend/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	router.Route("/auth", auth.Setup)
}

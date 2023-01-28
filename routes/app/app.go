package app

import (
	"node-backend/routes/app/manage"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/manage", manage.Setup)
	router.Post("/list", listApps)
}

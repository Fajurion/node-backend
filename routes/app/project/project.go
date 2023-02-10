package project

import (
	"node-backend/routes/app/project/containers"
	"node-backend/routes/app/project/events"
	"node-backend/routes/app/project/manage"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/fetch", fetch)
	router.Route("/manage", manage.Unauthorized)
	router.Route("/events", events.Unauthorized)
	router.Route("/containers", containers.Unauthorized)
}

func Authorized(router fiber.Router) {
}

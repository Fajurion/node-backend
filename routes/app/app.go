package app

import (
	"node-backend/routes/app/manage"
	"node-backend/routes/app/project"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/project", project.Unauthorized)
}

func Authorized(router fiber.Router) {
	router.Route("/manage", manage.Setup)
	router.Route("/project", project.Authorized)
	router.Post("/list", listApps)
}

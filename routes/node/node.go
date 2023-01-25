package node

import (
	"node-backend/routes/node/manage"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/init", manage.Setup)
}

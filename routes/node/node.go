package node

import (
	"node-backend/routes/node/connect"
	"node-backend/routes/node/manage"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/manage", manage.Setup)
	router.Post("/connect", connect.Connect)
	router.Post("/token", generateToken)
}

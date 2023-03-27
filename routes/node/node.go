package node

import (
	"node-backend/routes/node/connect"
	"node-backend/routes/node/manage"
	"node-backend/routes/node/status"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/status", status.Setup)
	router.Route("/manage", manage.Unauthorized)
	router.Post("/this", this)
	router.Post("/disconnect", connect.Disconnect)
}

func Authorized(router fiber.Router) {
	router.Route("/manage", manage.Authorized)
	router.Post("/connect", connect.Connect)
	router.Post("/token", generateToken)
}

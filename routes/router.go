package routes

import (
	"node-backend/routes/account/auth"
	"node-backend/routes/cluster"
	"node-backend/routes/node"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	router.Route("/auth", auth.Setup)
	router.Route("/node", node.Setup)
	router.Route("/cluster", cluster.Setup)
}

package routes

import (
	"node-backend/routes/account/auth"
	"node-backend/routes/cluster"
	"node-backend/routes/node"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Router(router fiber.Router) {

	// Unauthorized routes
	router.Route("/auth", auth.Setup)

	router.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(util.JWT_SECRET),
	}))

	// Authorized routes
	router.Route("/node", node.Setup)
	router.Route("/cluster", cluster.Setup)
}

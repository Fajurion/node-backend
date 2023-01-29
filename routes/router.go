package routes

import (
	"node-backend/routes/account/auth"
	"node-backend/routes/app"
	"node-backend/routes/cluster"
	"node-backend/routes/node"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Router(router fiber.Router) {

	// Unauthorized routes
	router.Route("/auth", auth.Setup)
	router.Route("/node", node.Unauthorized)

	router.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(util.JWT_SECRET),

		// Checks if the token is expired
		SuccessHandler: func(c *fiber.Ctx) error {

			if util.IsExpired(c) {
				return requests.InvalidRequest(c)
			}

			return c.Next()
		},
	}))

	// Authorized routes
	router.Route("/node", node.Authorized)
	router.Route("/app", app.Setup)
	router.Route("/cluster", cluster.Setup)
}

package routes

import (
	"log"
	"node-backend/routes/account"
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
	router.Route("/auth", auth.Unauthorized)
	router.Route("/node", node.Unauthorized)
	router.Route("/app", app.Unauthorized)
	router.Route("/account", account.Unauthorized)

	router.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(util.JWT_SECRET),

		// Checks if the token is expired
		SuccessHandler: func(c *fiber.Ctx) error {

			if util.IsExpired(c) {
				return requests.InvalidRequest(c)
			}

			return c.Next()
		},

		// Error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			log.Println(err.Error())

			// Return error message
			return c.SendStatus(401)
		},
	}))

	// Authorized routes
	router.Route("/account", account.Authorized)
	router.Route("/node", node.Authorized)
	router.Route("/app", app.Authorized)
	router.Route("/cluster", cluster.Setup)
}

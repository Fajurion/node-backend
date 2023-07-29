package routes

import (
	"log"
	nb_challenges "node-backend/nbchallenges"
	"node-backend/routes/account"
	"node-backend/routes/account/auth"
	"node-backend/routes/app"
	"node-backend/routes/cluster"
	"node-backend/routes/node"
	"node-backend/util"
	"node-backend/util/requests"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {

	// Unauthorized routes
	router.Route("/auth", auth.Unauthorized)
	router.Route("/node", node.Unauthorized)
	router.Route("/app", app.Unauthorized)
	router.Route("/account", account.Unauthorized)

	// Challenge test
	router.Post("/challenge/generate", nb_challenges.Generate)
	router.Post("/challenge/solve", nb_challenges.Solve)

	router.Route("/", remoteIDRoutes)
	router.Route("/", authorizedRoutes)
}

func remoteIDRoutes(router fiber.Router) {

	// Authorized by using a remote id or normal token
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(util.JWT_SECRET),
		},

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

	// Routes that require a remote id or normal JWT
	router.Route("/account", account.Remote)

}

func authorizedRoutes(router fiber.Router) {

	// Autorized by using a normal JWT token
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(util.JWT_SECRET),
		},

		// Checks if the token is expired
		SuccessHandler: func(c *fiber.Ctx) error {

			if util.IsExpired(c) {
				return requests.InvalidRequest(c)
			}

			if util.IsRemoteId(c) {
				requests.InvalidRequest(c)
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
	router.Route("/auth", auth.Authorized)
}

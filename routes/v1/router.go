package routes_v1

import (
	"log"
	nb_challenges "node-backend/nbchallenges"
	"node-backend/routes/v1/account"
	"node-backend/routes/v1/account/auth"
	"node-backend/routes/v1/app"
	"node-backend/routes/v1/cluster"
	"node-backend/routes/v1/node"
	"node-backend/util"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {

	// Get default private and public key
	serverPublicKey, err := util.UnpackageRSAPublicKey(os.Getenv("TC_PUBLIC_KEY"))
	if err != nil {
		panic("Couldn't unpackage public key. Required for v1 API. Please set TC_PUBLIC_KEY in your environment variables or .env file.")
	}

	serverPrivateKey, err := util.UnpackageRSAPrivateKey(os.Getenv("TC_PRIVATE_KEY"))
	if err != nil {
		panic("Couldn't unpackage private key. Required for v1 API. Please set TC_PRIVATE_KEY in your environment variables or .env file.")
	}

	// Through Cloudflare Protection (Decryption method)
	router.Use(func(c *fiber.Ctx) error {

		packagedPub, valid := c.GetReqHeaders()["Public-Key"]
		if !valid {
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}

		pub, err := util.UnpackageRSAPublicKey(packagedPub)
		if err != nil {
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}

		decrypted, err := util.DecryptRSA(serverPrivateKey, c.Body())
		if err != nil {
			return c.SendStatus(fiber.StatusNetworkAuthenticationRequired)
		}
		c.Locals("body", decrypted)
		c.Locals("pub", pub)
		c.Locals("srv_pub", serverPublicKey)

		return c.Next()
	})

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
				return util.InvalidRequest(c)
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
				return util.InvalidRequest(c)
			}

			if util.IsRemoteId(c) {
				util.InvalidRequest(c)
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

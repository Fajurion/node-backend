package routes_v1

import (
	"crypto/rsa"
	"encoding/base64"
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

	// Endpoint to get server public key (so no requirements apply yet)
	router.Post("/pub", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"pub": util.PackageRSAPublicKey(serverPublicKey),
		})
	})

	// Unencrypted account routes (only file upload thing)
	router.Route("/account", account.Unencrypted)

	router.Route("/", func(router fiber.Router) {
		encryptedRoutes(router, serverPublicKey, serverPrivateKey)
	})
}

func encryptedRoutes(router fiber.Router, serverPublicKey *rsa.PublicKey, serverPrivateKey *rsa.PrivateKey) {

	// Through Cloudflare Protection (Decryption method)
	router.Use(func(c *fiber.Ctx) error {

		aesKeyEncoded, valid := c.GetReqHeaders()["Auth-Tag"]
		if !valid {
			log.Println("no header")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}
		aesKeyEncrypted, err := base64.StdEncoding.DecodeString(aesKeyEncoded)
		if err != nil {
			log.Println("no decoding")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}

		aesKey, err := util.DecryptRSA(serverPrivateKey, aesKeyEncrypted)
		if err != nil {
			return c.SendStatus(fiber.StatusPreconditionRequired)
		}
		decrypted, err := util.DecryptAES(aesKey, c.Body())
		if err != nil {
			return c.SendStatus(fiber.StatusNetworkAuthenticationRequired)
		}
		c.Locals("body", decrypted)
		c.Locals("key", aesKey)
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

	router.Route("/", authorizedRoutes)
}

func authorizedRoutes(router fiber.Router) {

	// Autorized by using a normal JWT token
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS512,
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

	// Authorized routes
	router.Route("/account", account.Authorized)
	router.Route("/node", node.Authorized)
	router.Route("/app", app.Authorized)
	router.Route("/cluster", cluster.Setup)
	router.Route("/auth", auth.Authorized)
}

package routes_v1

import (
	"crypto/rsa"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

// Route: /pub
func getPublicKey(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"pub": util.PackageRSAPublicKey(c.Locals("srv_pub").(*rsa.PublicKey)),
	})
}

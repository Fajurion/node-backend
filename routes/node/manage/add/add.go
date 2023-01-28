package add

import "github.com/gofiber/fiber/v2"

func Setup(router fiber.Router) {
	router.Post("/new", newNode)
	router.Post("/existing", existingNode)
	router.Post("/token", createToken)
}

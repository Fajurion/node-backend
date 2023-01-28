package manage

import (
	"node-backend/routes/node/manage/add"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/add", add.Setup)
	router.Post("/remove", removeNode)
	router.Post("/regen", regenToken)
}

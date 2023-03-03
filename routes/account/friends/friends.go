package friends

import (
	"node-backend/routes/account/friends/request"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/request", request.Setup)
	router.Post("/check", checkFriendships)
}

func Authorized(router fiber.Router) {
	router.Post("/remove", removeFriend)
	router.Post("/list", listFriends)
	router.Post("/online", onlineFriends)
}

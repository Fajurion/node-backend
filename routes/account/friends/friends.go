package friends

import (
	"node-backend/routes/account/friends/request"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/request", request.Setup)
	router.Post("/check", checkFriendships)
	router.Post("/remove", removeFriend)
}

func Authorized(router fiber.Router) {
	router.Post("/list", listFriends)
	router.Post("/online", onlineFriends)
}

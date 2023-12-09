package account

import (
	"node-backend/routes/account/files"
	"node-backend/routes/account/friends"
	"node-backend/routes/account/keys"
	"node-backend/routes/account/profile"
	"node-backend/routes/account/rank"
	"node-backend/routes/account/stored_actions"
	"node-backend/routes/account/vault"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/rank", rank.Unauthorized)
	router.Route("/stored_actions", stored_actions.Unauthorized)
}

func Remote(router fiber.Router) {
	router.Route("/stored_actions", stored_actions.Remote)
	router.Route("/files", files.RemoteID)
	router.Route("/profile", profile.Remote)

	router.Post("/get", getAccount)
}

func Authorized(router fiber.Router) {
	router.Route("/keys", keys.Authorized)
	router.Route("/stored_actions", stored_actions.Authorized)
	router.Route("/friends", friends.Authorized)
	router.Route("/vault", vault.Authorized)
	router.Route("/profile", profile.Authorized)

	router.Post("/remote_id", generateRemoteId)
	router.Post("/me", me)
	router.Route("/files", files.Authorized)
}

package account

import (
	"node-backend/routes/v1/account/files"
	"node-backend/routes/v1/account/friends"
	invite_routes "node-backend/routes/v1/account/invite"
	"node-backend/routes/v1/account/keys"
	"node-backend/routes/v1/account/profile"
	"node-backend/routes/v1/account/rank"
	"node-backend/routes/v1/account/stored_actions"
	"node-backend/routes/v1/account/vault"

	"github.com/gofiber/fiber/v2"
)

func Unencrypted(router fiber.Router) {
	router.Route("/files", files.Unencrypted)
}

func Unauthorized(router fiber.Router) {
	router.Route("/rank", rank.Unauthorized)
	router.Route("/stored_actions", stored_actions.Unauthorized)
}

func Remote(router fiber.Router) {
	router.Route("/stored_actions", stored_actions.Remote)
	router.Route("/profile", profile.Remote)
	router.Route("/files", files.RemoteID)

	router.Post("/get", getAccount)
}

func Authorized(router fiber.Router) {
	router.Route("/keys", keys.Authorized)
	router.Route("/stored_actions", stored_actions.Authorized)
	router.Route("/friends", friends.Authorized)
	router.Route("/vault", vault.Authorized)
	router.Route("/profile", profile.Authorized)
	router.Route("/invite", invite_routes.Authorized)

	router.Post("/remote_id", generateRemoteId)
	router.Post("/me", me)
	router.Route("/files", files.Authorized)
}

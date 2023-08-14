package stored_actions

import "github.com/gofiber/fiber/v2"

// Configuration
const StoredActionLimit = 10       // Max number of stored actions per account
const StoredActionTokenLength = 32 // Length of the token used to identify stored actions

// Authorized with remote id
func Remote(router fiber.Router) {
	router.Post("/details", getDetails)
	router.Post("/send", sendStoredAction)
}

// Authorized with account JWT
func Authorized(router fiber.Router) {
	router.Post("/list", listStoredActions)
	router.Post("/delete", deleteStoredAction)
}

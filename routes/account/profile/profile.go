package profile

import "github.com/gofiber/fiber/v2"

// Route manager for /account/profile (authorized)
func SetupRoutes(router fiber.Router) {

	router.Post("/disable", disableProfile)
	router.Post("/enable", enableProfile)
	router.Post("/edit", editProfile)
	router.Post("/get", getProfile)
	router.Post("/me", currentProfile)
}

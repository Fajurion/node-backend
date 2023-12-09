package profile

import "github.com/gofiber/fiber/v2"

func Authorized(router fiber.Router) {
	router.Post("/set_picture", setProfilePicture)
}

func Remote(router fiber.Router) {

}

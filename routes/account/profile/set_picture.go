package profile

import (
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type setProfileRequest struct {
	File string `json:"file"`
}

var fileTypes = []string{
	"png",
	"jpg",
	"jpeg",
}

// Route: /account/profile/set_picture
func setProfilePicture(c *fiber.Ctx) error {

	var req setProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	return requests.SuccessfulRequest(c)
}

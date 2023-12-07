package files

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type infoRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/info
func fileInfo(c *fiber.Ctx) error {

	var req infoRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get file info
	var cloudFile account.CloudFile
	if err := database.DBConn.Select("id,name,size,account").Where("id = ?", req.Id).Take(&cloudFile).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"file":    cloudFile,
	})
}

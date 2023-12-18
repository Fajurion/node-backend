package files

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type unfavoriteRequest struct {
	Id string `json:"id"`
}

func unfavoriteFile(c *fiber.Ctx) error {

	var req unfavoriteRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}
	accId := util.GetAcc(c)

	// Get file
	var file account.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return requests.FailedRequest(c, "file.not_found", nil)
	}

	// Toggle favorite
	if err := database.DBConn.Model(&account.CloudFile{}).Update("favorite", false).Where("id = ?", file.Id).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

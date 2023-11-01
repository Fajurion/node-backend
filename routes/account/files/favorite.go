package files

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type favoriteRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/favorite
func favoriteFile(c *fiber.Ctx) error {

	var req favoriteRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}
	accId := util.GetAcc(c)

	// Get file
	var file account.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return requests.FailedRequest(c, "file.not_found", nil)
	}

	favoriteStorage, err := CountFavoriteStorage(accId)
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	if favoriteStorage+file.Size > maxFavoriteStorage {
		return requests.FailedRequest(c, "file.favorite_limit", nil)
	}

	// Toggle favorite
	if err := database.DBConn.Model(&account.CloudFile{}).Update("favorite", true).Where("id = ?", file.Id).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

package files

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type favoriteRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/favorite
func favoriteFile(c *fiber.Ctx) error {

	var req favoriteRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId := util.GetAcc(c)

	// Get file
	var file account.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return util.FailedRequest(c, "file.not_found", nil)
	}

	favoriteStorage, err := CountFavoriteStorage(accId)
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	if favoriteStorage+file.Size > maxFavoriteStorage {
		return util.FailedRequest(c, "file.favorite_limit", nil)
	}

	// Toggle favorite
	if err := database.DBConn.Model(&account.CloudFile{}).Update("favorite", true).Where("id = ?", file.Id).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

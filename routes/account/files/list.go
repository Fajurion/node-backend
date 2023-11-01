package files

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type listRequest struct {
	Favorite bool  `json:"favorite"`
	Start    int64 `json:"last"` // Start data
}

// Route: /account/files/list
func listFiles(c *fiber.Ctx) error {

	var req listRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	accId := util.GetAcc(c)

	// Get files
	var files []account.CloudFile
	if database.DBConn.Where("account = ? AND favorite = ? AND created_at < ?", accId, req.Start).Limit(40).Find(&[]account.CloudFile{}).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"file":    files,
	})
}

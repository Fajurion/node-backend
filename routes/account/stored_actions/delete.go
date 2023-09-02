package stored_actions

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type deleteRequest struct {
	ID string `json:"id"`
}

// Route: /account/stored_actions/delete
func deleteStoredAction(c *fiber.Ctx) error {

	// Parse request
	var req deleteRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Delete stored action
	accId := util.GetAcc(c)
	if err := database.DBConn.Where("account = ? AND id = ?", accId, req.ID).Delete(&properties.StoredAction{}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}
	if err := database.DBConn.Where("account = ? AND id = ?", accId, req.ID).Delete(&properties.AStoredAction{}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Return success
	return requests.SuccessfulRequest(c)
}

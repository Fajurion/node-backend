package stored_actions

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type sendRequest struct {
	Account string `json:"account"`
	Payload string `json:"payload"`
}

func sendStoredAction(c *fiber.Ctx) error {

	// Parse request
	var req sendRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where("id = ?", req.Account).Take(&acc).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if stored action limit is reached
	var storedActionCount int64
	if err := database.DBConn.Model(&properties.StoredAction{}).Where("account = ?", acc.ID).Count(&storedActionCount).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	if storedActionCount >= StoredActionLimit {
		return requests.FailedRequest(c, "limit.reached", nil)
	}

	// Create stored action
	storedAction := properties.StoredAction{
		ID:      auth.GenerateToken(12),
		Account: acc.ID,
		Payload: req.Payload,
	}

	// Save stored action
	if err := database.DBConn.Create(&storedAction).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Send to node if possible
	var session account.Session
	if err := database.DBConn.Where("account = ? AND node != ?", acc.ID, 0).Take(&session).Error; err == nil {

		// No error handling, cause it doesn't matter if it couldn't send
		requests.SendEventToNode(session.Node, acc.ID, requests.Event{
			Sender: "0",
			Name:   "s_a", // Stored action
			Data: map[string]interface{}{
				"id":      storedAction.ID,
				"payload": storedAction.Payload,
			},
		})
	}

	// Return success
	return requests.SuccessfulRequest(c)
}

package stored_actions

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/stored_actions/send_auth
func sendAuthenticatedStoredAction(c *fiber.Ctx) error {

	// Parse request
	var req sendRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if req.Account == "" || req.Payload == "" {
		return requests.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where("id = ?", req.Account).Take(&acc).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	// Create stored action
	storedAction := properties.StoredAction{
		ID:      auth.GenerateToken(12),
		Account: acc.ID,
		Payload: req.Payload,
	}

	// Check if stored action limit is reached
	var storedActionCount int64
	if err := database.DBConn.Model(&properties.AStoredAction{}).Where("account = ?", acc.ID).Count(&storedActionCount).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	if storedActionCount >= AuthenticatedStoredActionLimit {
		return requests.FailedRequest(c, "limit.reached", nil)
	}

	var storedActionKey account.StoredActionKey
	if err := database.DBConn.Where(&account.StoredActionKey{ID: req.Account}).Take(&storedActionKey).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	if storedActionKey.Key != req.Key {
		return requests.InvalidRequest(c)
	}

	// Save authenticated stored action
	if err := database.DBConn.Create(&properties.AStoredAction{
		ID:      storedAction.ID,
		Account: storedAction.Account,
		Payload: storedAction.Payload,
	}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Send stored action to account
	sendStoredActionTo(acc.ID, true, storedAction)

	return requests.SuccessfulRequest(c)
}

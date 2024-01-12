package settings_routes

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

type changePasswordRequest struct {
	Current string `json:"current"`
	New     string `json:"new"`
}

// Change the password of an account (Route: /account/settings/change_password)
func changePassword(c *fiber.Ctx) error {

	var req changePasswordRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get current password
	accId := util.GetAcc(c)
	var authentication account.Authentication
	if err := database.DBConn.Where("account = ? AND type = ?", accId, account.TypePassword).Take(&authentication).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	log.Println(auth.HashPassword(req.Current, accId))

	// Check password
	if auth.HashPassword(req.Current, accId) != authentication.Secret {
		return util.FailedRequest(c, util.PasswordInvalid, nil)
	}

	// Log out all devices
	// TODO: Disconnect all sessions
	if err := database.DBConn.Where("account = ?", accId).Delete(&account.Session{}).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Change password
	err := database.DBConn.Model(&account.Authentication{}).Where("account = ? AND type = ?", accId, account.TypePassword).
		Update("secret", auth.HashPassword(req.New, accId)).Error
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// TODO: Send a mail here in the future (Stuff required: Rate limiting)

	return util.SuccessfulRequest(c)
}

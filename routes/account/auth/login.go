package auth

import (
	"errors"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoginRequest is the request body for the login request
type startLoginRequest struct {
	Email  string `json:"email"`
	Device string `json:"device"`
}

// startLogin starts the login process
func startLogin(c *fiber.Ctx) error {

	var req startLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if user exists
	var acc account.Account
	if database.DBConn.Where("email = ?", req.Email).Take(&acc).Error != nil {
		return requests.InvalidRequest(c)
	}

	// Generate token
	tk, err := auth.GenerateLoginTokenWithStep(acc.ID, req.Device, 1)
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Get authentication methods for step
	var methods []uint
	if database.DBConn.Where("step = ? AND account = ?", account.StartStep, acc.ID).Select("type").Take(&methods).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   tk,
		"methods": methods,
	})
}

type loginStepRequest struct {
	Secret string `json:"secret"`
}

// loginStep runs the login step
func loginStep(c *fiber.Ctx) error {

	// TODO: Auth stuff

	// Create session
	tk := auth.GenerateToken(100)

	var createdSession account.Session = account.Session{
		ID:              auth.GenerateToken(8),
		Token:           tk,
		Account:         "23",
		PermissionLevel: 0,
		Device:          "ph", // TODO: Give the user the option to choose the device
		LastConnection:  time.UnixMilli(0),
	}

	if err := database.DBConn.Create(&createdSession).Error; err != nil {
		requests.FailedRequest(c, "server.error", err)
	}

	// Generate jwt token
	jwtToken, err := util.Token(createdSession.ID, "23", 0, time.Now().Add(time.Hour*24*3))

	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"token":         jwtToken,
		"refresh_token": tk,
	})
}

// checkSessions checks if the user has too many sessions
func checkSessions(c *fiber.Ctx, acc account.Account) (bool, error) {

	// Check if user has too many sessions
	var sessions int64
	if err := database.DBConn.Model(&account.Session{}).Where("account = ?", acc.ID).Count(&sessions).Error; err != nil {
		return false, errors.New("server.error")
	}

	if sessions > 10 {
		return false, nil
	}

	return true, nil
}

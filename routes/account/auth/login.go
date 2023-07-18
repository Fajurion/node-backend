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

// startLogin starts the login process Route: /auth/login/start
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

	valid, err := checkSessions(acc.ID)
	if err != nil {
		return requests.FailedRequest(c, err.Error(), nil)
	}

	if !valid {
		return requests.FailedRequest(c, "too.many.sessions", nil)
	}

	// Generate token
	return runAuthStep(acc.ID, req.Device, account.StartStep, c)
}

type loginStepRequest struct {
	Type   uint   `json:"type"`
	Secret string `json:"secret"`
}

// loginStep runs the login step Route: /auth/login/step
func loginStep(c *fiber.Ctx) error {

	var req loginStepRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get data
	id, device, step := auth.GetLoginDataFromToken(c)

	var method account.Authentication
	if err := database.DBConn.Where("account = ? AND type = ?", id, req.Type).Take(&method).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Check the provided secret
	if !method.Verify(req.Type, req.Secret) {
		return requests.FailedRequest(c, "invalid.method", nil)
	}

	return runAuthStep(id, device, step+1, c)
}

func runAuthStep(id string, device string, step uint, c *fiber.Ctx) error {

	// Generate token
	tk, err := auth.GenerateLoginTokenWithStep(id, device, step)
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Get authentication methods for step
	var availableMethods []uint
	for method, stepUsed := range account.Order {
		if stepUsed == step {
			availableMethods = append(availableMethods, method)
		}
	}

	var methods []uint
	if database.DBConn.Where("type IN ? AND account = ?", availableMethods, id).Select("type").Take(&methods).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	if len(methods) == 0 && step == account.StartStep {
		// TODO: SERIOUS SECURITY ISSUE WARNING HERE
		return requests.FailedRequest(c, "no.methods", nil)
	}

	if len(methods) == 0 {

		var acc account.Account
		if err := database.DBConn.Where("id = ?", id).Preload("Rank").Take(&acc).Error; err != nil {
			return requests.FailedRequest(c, "server.error", err)
		}

		// Create session
		tk := auth.GenerateToken(100)

		var createdSession account.Session = account.Session{
			ID:              auth.GenerateToken(8),
			Token:           tk,
			Account:         acc.ID,
			PermissionLevel: acc.Rank.Level,
			Device:          device,
			LastConnection:  time.UnixMilli(0),
		}

		if err := database.DBConn.Create(&createdSession).Error; err != nil {
			requests.FailedRequest(c, "server.error", err)
		}

		// Generate jwt token
		jwtToken, err := util.Token(createdSession.ID, acc.ID, acc.Rank.Level, time.Now().Add(time.Hour*24*3))

		if err != nil {
			return requests.FailedRequest(c, "server.error", err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":       true,
			"token":         jwtToken,
			"refresh_token": tk,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   tk,
		"methods": methods,
	})
}

// checkSessions checks if the user has too many sessions
func checkSessions(id string) (bool, error) {

	// Check if user has too many sessions
	var sessions int64
	if err := database.DBConn.Model(&account.Session{}).Where("account = ?", id).Count(&sessions).Error; err != nil {
		return false, errors.New("server.error")
	}

	if sessions > 10 {
		return false, nil
	}

	return true, nil
}

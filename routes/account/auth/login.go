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
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// login handles the login request
func login(c *fiber.Ctx) error {
	var req loginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if request is valid
	if req.Email == "" || req.Password == "" {
		return requests.FailedRequest(c, "invalid", nil)
	}

	// Get user from database
	var acc account.Account
	if err := database.DBConn.Where("email = ?", req.Email).Preload("Rank").First(&acc).Error; err != nil {
		return requests.FailedRequest(c, "invalid.password", nil)
	}

	// Check account details
	if !checkAccountDetails(c, acc, req) {
		return requests.FailedRequest(c, "invalid.password", nil)
	}

	// Check if user has too many sessions
	if valid, err := checkSessions(c, acc); err != nil || !valid {
		return requests.FailedRequest(c, "too.many.sessions", nil)
	}

	// Create session
	tk := auth.GenerateToken(100)

	var createdSession account.Session = account.Session{
		Token:           tk,
		Account:         acc.ID,
		PermissionLevel: acc.Rank.Level,
		Device:          "ph", // TODO: Give the user the option to choose the device
		End:             time.Now().Add(time.Hour * 24 * 7),
		LastConnection:  time.UnixMilli(0),
	}

	if err := database.DBConn.Create(&createdSession).Error; err != nil {
		requests.FailedRequest(c, "server.error", err)
	}

	// Generate jwt token
	jwtToken, err := util.Token(createdSession.ID, time.Now().Add(time.Hour*24*3), fiber.Map{
		"acc": acc.ID,
		"lvl": acc.Rank.Level,
	})

	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"token":         jwtToken,
		"refresh_token": tk,
	})
}

// checkAccountDetails checks if the account details are valid
func checkAccountDetails(c *fiber.Ctx, acc account.Account, req loginRequest) bool {

	if acc.ID == 0 {
		return false
	}

	if !acc.CheckPassword(req.Password) {
		return false
	}

	return true
}

// checkSessions checks if the user has too many sessions
func checkSessions(c *fiber.Ctx, acc account.Account) (bool, error) {

	// Check if user has too many sessions
	var sessions int64
	if err := database.DBConn.Where(&account.Session{Account: acc.ID}).Count(&sessions).Error; err != nil {
		return false, errors.New("server.error")
	}

	if sessions > 10 {
		return false, nil
	}

	return true, nil
}

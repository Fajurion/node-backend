package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Post("/login", login)
	router.Post("/register", register_test)
}

// LoginRequest is the request body for the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// login handles the login request
func login(c *fiber.Ctx) error {
	var req LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid",
		})
	}

	// Check if request is valid
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid",
		})
	}

	// Get user from database
	var acc account.Account
	database.DBConn.Model(&account.Account{}).Where("email = ?", req.Email).Preload("Rank").First(&acc)

	// Check account details
	if err := checkAccountDetails(c, acc, req); err != nil {
		return err
	}

	// Check if user has too many sessions
	if err := checkSessions(c, acc); err != nil {
		return err
	}

	// Create session
	token := auth.GenerateToken()

	err := database.DBConn.Create(&account.Session{
		Token:           token,
		Account:         acc.ID,
		PermissionLevel: acc.Rank.Level,
		Device:          "web", // TODO: Get device from request
		Connected:       false,
	}).Error

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "server.error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"token":   token,
		"message": "success",
	})
}

// checkAccountDetails checks if the account details are valid
func checkAccountDetails(c *fiber.Ctx, acc account.Account, req LoginRequest) error {

	if acc.ID == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid.password",
		})
	}

	if !acc.CheckPassword(acc.Password) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid.password",
		})
	}

	return nil
}

// checkSessions checks if the user has too many sessions
func checkSessions(c *fiber.Ctx, acc account.Account) error {

	// Check if user has too many sessions
	var sessions []account.Session
	database.DBConn.Where("account = ?", acc.ID).Find(&sessions)

	// TODO: Max sessions in application properties
	if len(sessions) > 10 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "too.many.sessions",
		})
	}

	connected := 0
	for _, session := range sessions {
		if session.Connected {
			connected++
		}
	}

	// TODO: Max connected sessions in application properties
	if connected > 3 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "too.many.connected",
		})
	}

	return nil
}

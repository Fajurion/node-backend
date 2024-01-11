package auth

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/standards"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Request sent to start the registration process
type registerStartRequest struct {
	Locale string `json:"locale"`
	Email  string `json:"email"`
	Invite string `json:"invite"`
}

// Claims in the JWT for email verification (step 2, generated by step 1)
type registerEmailTokenClaims struct {
	Step           int    `json:"s"`
	Email          string `json:"e"`
	Code           string `json:"c"`
	ExpiredUnixSec int64  `json:"e_u"`

	jwt.RegisteredClaims
}

// Route: /auth/register/start, start the registration process
func registerStart(c *fiber.Ctx) error {

	// Parse body to register request
	var req registerStartRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check the invite
	var invite account.Invite
	if err := database.DBConn.Where("id = ?", req.Invite).Take(&invite).Error; err != nil {
		return util.FailedRequest(c, util.InviteInvalid, err)
	}

	// Just for security
	if invite.ID != req.Invite {
		return util.FailedRequest(c, util.InviteInvalid, nil)
	}

	// Check if email is already registered
	valid, normalizedEmail := standards.CheckEmail(req.Email)
	if !valid {
		return util.FailedRequest(c, util.EmailInvalid, nil)
	}

	if database.DBConn.Where("email = ?", normalizedEmail).Take(&account.Account{}).RowsAffected > 0 {
		return util.FailedRequest(c, util.EmailRegistered, nil)
	}

	// Generate a hidden jwt value for the verification code of the email
	code := auth.GenerateToken(6)
	codeHV, err := util.MakeHiddenJWTValue(c, []byte(code)) // HV = hidden value
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Generate a start registration token
	exp := time.Now().Add(time.Hour * 2)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, registerEmailTokenClaims{
		Step:           1,
		Email:          normalizedEmail,
		Code:           codeHV,
		ExpiredUnixSec: exp.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(util.JWT_SECRET))
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Send email
	err = mail.SendEmail(normalizedEmail, req.Locale, mail.EmailVerification, code)
	if err != nil {
		return util.FailedRequest(c, util.ErrorMail, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   tokenString,
	})
}

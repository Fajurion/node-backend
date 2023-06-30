package auth

import (
	"node-backend/util"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateLoginTokenWithStep(id string, device string, step uint) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"s":   step,
		"e_u": time.Now().Add(time.Minute * 5).Unix(), // Expiration unix
		"acc": id,
		"d":   device,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(util.JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

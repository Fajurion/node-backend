package util

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v4"
)

func Token(token string, exp time.Time, data map[string]interface{}) (string, error) {

	jsonData, err := sonic.Marshal(data)

	if err != nil {
		return "", err
	}

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tk":   token,
		"exp":  exp.Unix(),
		"data": string(jsonData),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

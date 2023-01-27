package util

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWT_SECRET is the secret used to sign the jwt token
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

// IsExpired checks if the token is expired
func IsExpired(c *fiber.Ctx) bool {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	exp := claims["exp"].(int64)

	return time.Now().Unix() > exp
}

// Permission checks if the user has the required permission level
func Permission(c *fiber.Ctx, perm int16) bool {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})

	return data["lvl"].(int16) >= perm
}

func GetData(c *fiber.Ctx) map[string]interface{} {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})

	return data
}

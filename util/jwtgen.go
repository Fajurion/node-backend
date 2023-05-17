package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWT_SECRET is the secret used to sign the jwt token
func Token(session string, exp time.Time, data map[string]interface{}) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ses":  session,
		"e_u":  exp.Unix(), // Expiration unix
		"data": data,
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

	num := claims["e_u"].(float64)
	exp := int64(num)

	return time.Now().Unix() > exp
}

// Permission checks if the user has the required permission level
func Permission(c *fiber.Ctx, perm int16) bool {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})

	num := data["lvl"].(float64)
	lvl := int16(num)

	return lvl >= perm
}

func GetSession(c *fiber.Ctx) string {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["ses"].(string)
}

func GetData(c *fiber.Ctx) map[string]interface{} {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})

	return data
}

package util

import (
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWT_SECRET is the secret used to sign the jwt token
func Token(session string, account string, lvl uint, exp time.Time) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ses": session,
		"e_u": exp.Unix(), // Expiration unix
		"acc": account,
		"lvl": lvl,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RemoteId(lvl uint, random string) (string, error) {

	// Create jwt token
	exp := time.Now().Add(time.Hour * 2)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"e_u": exp.Unix(), // Expiration unix
		"lvl": lvl,
		"r":   random,
		"rid": true, // tell the backend that it's a remote id
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
func Permission(c *fiber.Ctx, perm string) bool {

	// Check if there is a JWT token
	if c.Locals("user") == nil || reflect.TypeOf(c.Locals("user")).String() != "*jwt.Token" {
		return false
	}

	// Get the permission from the map
	permission, valid := Permissions[perm]
	if !valid {
		return false
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	lvl := int16(claims["lvl"].(float64))

	return lvl >= permission
}

func GetSession(c *fiber.Ctx) string {
	if c.Locals("user") == nil || reflect.TypeOf(c.Locals("user")).String() != "*jwt.Token" {
		return ""
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["ses"].(string)
}

func GetAcc(c *fiber.Ctx) string {
	if c.Locals("user") == nil || reflect.TypeOf(c.Locals("user")).String() != "*jwt.Token" {
		return ""
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["acc"].(string)
}

func IsRemoteId(c *fiber.Ctx) bool {
	if c.Locals("user") == nil || reflect.TypeOf(c.Locals("user")).String() != "*jwt.Token" {
		return false
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	_, ok := claims["rid"]
	return ok
}

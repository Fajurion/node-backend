package auth

import (
	"crypto/sha256"
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateToken() string {

	// Generate a 200 characters long random string
	s := make([]rune, 200)

	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

// HashPassword hashes a password
func HashPassword(password string) string {

	// Get hash of password
	hash := sha256.Sum256([]byte(password))

	// Convert byte[] to string
	return string(hash[:])
}

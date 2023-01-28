package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateToken(tkLength int32) string {

	s := make([]rune, tkLength)

	length := big.NewInt(int64(len(letters)))

	for i := range s {

		number, _ := rand.Int(rand.Reader, length)
		s[i] = letters[number.Int64()]
	}

	return string(s)
}

// HashPassword hashes a password
func HashPassword(password string) string {

	// Get hash of password
	hash := sha256.Sum256([]byte(password))

	// Convert byte[] to string
	return base64.StdEncoding.EncodeToString(hash[:])
}

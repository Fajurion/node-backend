package util

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"testing"
)

func TestDecryptRSA(t *testing.T) {
	// Generate a new RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	// Encrypt a sample plaintext using the public key
	plaintext := "Hello, World!"
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, []byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt plaintext: %v", err)
	}

	// Convert the ciphertext to a base64-encoded string
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)

	// Decrypt the ciphertext using the private key
	decryptedPlaintext, err := DecryptRSA(privateKey, ciphertextBase64)
	if err != nil {
		t.Fatalf("Failed to decrypt ciphertext: %v", err)
	}

	// Check if the decrypted plaintext matches the original plaintext
	if decryptedPlaintext != plaintext {
		t.Errorf("Decrypted plaintext does not match original plaintext. Expected: %s, Got: %s", plaintext, decryptedPlaintext)
	}
}

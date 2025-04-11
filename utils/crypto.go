package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/MiguelMachado-dev/disc-go-bot/config"
)

var cryptoLog = config.NewLogger("crypto")

// Encrypt encrypts data using AES-GCM
func Encrypt(plaintext string) (string, error) {
	key := []byte(config.GetEnv().ENCRYPTION_KEY)

	// Ensure key is exactly 32 bytes (pad or truncate)
	key = padKey(key, 32)

	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return base64 encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-GCM
func Decrypt(encryptedText string) (string, error) {
	key := []byte(config.GetEnv().ENCRYPTION_KEY)

	// Ensure key is exactly 32 bytes (pad or truncate)
	key = padKey(key, 32)

	// Decode base64 string
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract the nonce and actual ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// padKey ensures the key is exactly the required length by either padding or truncating
func padKey(key []byte, requiredLength int) []byte {
	currentLength := len(key)

	// If key is already the right length, return it
	if currentLength == requiredLength {
		return key
	}

	// If key is too long, truncate it
	if currentLength > requiredLength {
		return key[:requiredLength]
	}

	// If key is too short, pad it with zeros
	paddedKey := make([]byte, requiredLength)
	copy(paddedKey, key)
	return paddedKey
}

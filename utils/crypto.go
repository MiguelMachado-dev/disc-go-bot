package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/MiguelMachado-dev/disc-go-bot/config"
)

var cryptoLog = config.NewLogger("crypto")

// Encrypt encrypts data using AES-GCM
func Encrypt(plaintext string) (string, error) {
	key := []byte(config.GetEnv().ENCRYPTION_KEY)

	// Ensure key is exactly 32 bytes
	if len(key) != 32 {
		errMsg := fmt.Sprintf("encryption key must be exactly 32 bytes long, got %d bytes", len(key))
		cryptoLog.Errorf(errMsg)
		return "", errors.New(errMsg)
	}

	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		cryptoLog.Errorf("error creating cipher block: %v", err)
		return "", err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		cryptoLog.Errorf("error creating GCM: %v", err)
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		cryptoLog.Errorf("error creating nonce: %v", err)
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

	// Ensure key is exactly 32 bytes
	if len(key) != 32 {
		errMsg := fmt.Sprintf("encryption key must be exactly 32 bytes long, got %d bytes", len(key))
		cryptoLog.Errorf(errMsg)
		return "", errors.New(errMsg)
	}

	// Decode base64 string
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		cryptoLog.Errorf("error decoding base64 string: %v", err)
		return "", err
	}

	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		cryptoLog.Errorf("error creating cipher block: %v", err)
		return "", err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		cryptoLog.Errorf("error creating GCM: %v", err)
		return "", err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		errMsg := "ciphertext too short"
		cryptoLog.Errorf(errMsg)
		return "", errors.New(errMsg)
	}

	// Extract the nonce and actual ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		cryptoLog.Errorf("error decrypting data: %v", err)
		return "", err
	}

	return string(plaintext), nil
}

package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

func GenerateEncryptedToken(stringID string) (string, error) {
	encryptionKey, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	// Create a GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the subscriptionID
	ciphertext := gcm.Seal(nonce, nonce, []byte(stringID), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptToken(encryptedToken string) (string, error) {
	encryptionKey, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func getEncryptionKey() ([]byte, error) {
	key := os.Getenv("NEWSLETTER_ENCRYPTION_KEY")
	if len(key) != 32 { // AES-256 requires 32 bytes
		return nil, fmt.Errorf("invalid key length")
	}
	return []byte(key), nil
}

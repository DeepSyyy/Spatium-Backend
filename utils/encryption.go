package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// Encrypt text using AES-GCM and return base64 string
func Encrypt(text string) (string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	if len(key) < 32 {
		return "", errors.New("encryption key must be at least 32 bytes")
	}

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt base64-encoded AES-GCM ciphertext
func Decrypt(encrypted string) (string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	if len(key) < 32 {
		return "", errors.New("encryption key must be at least 32 bytes")
	}

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(data) < aesGCM.NonceSize() {
		return "", errors.New("invalid ciphertext length")
	}

	nonce, ciphertext := data[:aesGCM.NonceSize()], data[aesGCM.NonceSize():]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Shield provides AES-GCM encryption with Argon2id key derivation.
type Shield struct {
	passphrase string
}

// NewShield creates a new Shield with the given passphrase.
func NewShield(passphrase string) *Shield {
	return &Shield{
		passphrase: passphrase,
	}
}

const (
	saltSize   = 16
	nonceSize  = 12 // Default for AES-GCM
	keySize    = 32 // AES-256
	iterations = 3
	memory     = 64 * 1024
	parallel   = 4
)

// Encrypt encrypts the plaintext using Argon2id for key derivation.
// The output is a string in the format "salt.nonce.ciphertext" encoded in Base64.
func (s *Shield) Encrypt(plaintext string) (string, error) {
	// 1. Generate random salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	// 2. Derive key using Argon2id
	key := argon2.IDKey([]byte(s.passphrase), salt, iterations, memory, parallel, keySize)

	// 3. Initialize AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 4. Generate unique nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 5. Encrypt
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// 6. Bundle as salt.nonce.ciphertext
	output := fmt.Sprintf("%s.%s.%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(nonce),
		base64.RawStdEncoding.EncodeToString(ciphertext),
	)

	return output, nil
}

// Decrypt decrypts the base64 encoded string (which includes salt and nonce).
func (s *Shield) Decrypt(encoded string) (string, error) {
	parts := strings.Split(encoded, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid encrypted format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	nonce, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return "", err
	}

	// 1. Re-derive key using preserved salt
	key := argon2.IDKey([]byte(s.passphrase), salt, iterations, memory, parallel, keySize)

	// 2. Initialize AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. Decrypt
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}

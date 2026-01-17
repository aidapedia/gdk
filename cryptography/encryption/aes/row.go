package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// RowEncryptor encrypts/decrypts multiple fields with a single derived DEK.
// Use DeriveRow to create one. Always defer Close() to zero the DEK from memory.
type RowEncryptor struct {
	dek []byte
	gcm cipher.AEAD
	aad []byte
}

// EncryptField encrypts a single plaintext field value and returns base64-encoded ciphertext.
// Each call generates a unique nonce, so encrypting the same value twice yields different ciphertexts.
func (r *RowEncryptor) EncryptField(plaintext []byte) (string, error) {
	nonce := make([]byte, r.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := r.gcm.Seal(nonce, nonce, plaintext, r.aad)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptField decrypts a base64-encoded ciphertext field value and returns the plaintext string.
func (r *RowEncryptor) DecryptField(ciphertextBase64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}
	nonceSize := r.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", io.ErrUnexpectedEOF
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := r.gcm.Open(nil, nonce, ciphertext, r.aad)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// Close zeros the plaintext DEK from memory. Always defer this after DeriveRow.
func (r *RowEncryptor) Close() {
	zeroDEK(r.dek)
}

// newRowEncryptor creates a RowEncryptor from a plaintext DEK.
func newRowEncryptor(dek, aad []byte) (*RowEncryptor, error) {
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &RowEncryptor{
		dek: dek,
		gcm: gcm,
		aad: aad,
	}, nil
}

package dek

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Shield represents a DEK shield that encrypts and decrypts data using AES.
type Shield struct {
	kek []byte
}

// NewShield creates a new Shield instance with the provided KEK.
func NewShield(kek []byte) *Shield {
	return &Shield{
		kek: kek,
	}
}

// encrypt encrypts the DEK using the provided KEK.
func (s *Shield) encrypt(plaintextDEK []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.kek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintextDEK, nil), nil
}

// decrypt decrypts the DEK using the provided KEK.
func (s *Shield) decrypt(wrappedDEK []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.kek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(wrappedDEK) < nonceSize {
		return nil, io.ErrUnexpectedEOF
	}
	nonce, ciphertext := wrappedDEK[:nonceSize], wrappedDEK[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncryptedRecord holds the components you must save to your database.
type EncryptedRecord struct {
	Ciphertext string // The encrypted PII
	WrappedDEK string // The DEK, encrypted by your Master Key (KEK)
}

func (s *Shield) EnvelopeEncrypt(pii []byte) (*EncryptedRecord, error) {
	// 1. Generate a random 32-byte DEK locally
	dek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, err
	}

	// 2. Encrypt the PII using this local DEK
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nonce, nonce, pii, nil)

	// 3. Wrap the DEK using the remote Master Key (KEK)
	wrappedDEK, err := s.encrypt(dek)
	if err != nil {
		return nil, err
	}

	return &EncryptedRecord{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		WrappedDEK: base64.StdEncoding.EncodeToString(wrappedDEK),
	}, nil
}

func (s *Shield) EnvelopeDecrypt(record *EncryptedRecord) (string, error) {
	// 1. Decode base64 strings
	wrappedDEK, err := base64.StdEncoding.DecodeString(record.WrappedDEK)
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(record.Ciphertext)
	if err != nil {
		return "", err
	}

	// 2. Unwrap the DEK using the remote KMS
	dek, err := s.decrypt(wrappedDEK)
	if err != nil {
		return "", err
	}

	// 3. Decrypt the PII using the recovered DEK
	block, err := aes.NewCipher(dek)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return "", io.ErrUnexpectedEOF
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	return string(plaintext), err
}

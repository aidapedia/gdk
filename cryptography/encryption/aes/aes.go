package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	gcrypt "github.com/aidapedia/gdk/cryptography"
	gencryption "github.com/aidapedia/gdk/cryptography/encryption"
	"golang.org/x/crypto/hkdf"
)

const algorithmAES256GCM = "AES-256-GCM"

// AES represents a DEK AES that encrypts and decrypts data using AES.
type AES struct {
	kek         []byte
	dekByteSize int
	kekVersion  int
}

// NewAES creates a new AES instance with the provided KEK.
// kekVersion identifies which KEK is being used (for key rotation).
func NewAES(kek []byte, dekByteSize int, kekVersion int) (gcrypt.EncryptionInterface, error) {
	switch len(kek) {
	case 16, 24, 32:
	default:
		return nil, fmt.Errorf("invalid KEK size: %d bytes (must be 16, 24, or 32)", len(kek))
	}
	switch dekByteSize {
	case 16, 24, 32:
	default:
		return nil, fmt.Errorf("invalid DEK size: %d bytes (must be 16, 24, or 32)", dekByteSize)
	}
	return &AES{
		kek:         kek,
		dekByteSize: dekByteSize,
		kekVersion:  kekVersion,
	}, nil
}

// encrypt encrypts the DEK using the provided KEK.
func (s *AES) Encrypt(plaintextDEK []byte) ([]byte, error) {
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
func (s *AES) Decrypt(wrappedDEK []byte) ([]byte, error) {
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

func (s *AES) EncryptRecord(pii []byte, aad []byte) (*gencryption.EncryptedRecord, error) {
	// 1. Generate a DEK locally
	dek := make([]byte, s.dekByteSize)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, err
	}
	defer zeroDEK(dek)

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
	ciphertext := aesgcm.Seal(nonce, nonce, pii, aad)

	// 3. Wrap the DEK using the remote Master Key (KEK)
	wrappedDEK, err := s.Encrypt(dek)
	if err != nil {
		return nil, err
	}

	return &gencryption.EncryptedRecord{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		WrappedDEK: base64.StdEncoding.EncodeToString(wrappedDEK),
		KEKVersion: s.kekVersion,
		Algorithm:  algorithmAES256GCM,
	}, nil
}

func (s *AES) DecryptRecord(record *gencryption.EncryptedRecord, aad []byte) (string, error) {
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
	dek, err := s.Decrypt(wrappedDEK)
	if err != nil {
		return "", err
	}
	defer zeroDEK(dek)

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

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, aad)
	return string(plaintext), err
}

// zeroDEK overwrites the DEK slice with zeros.
func zeroDEK(dek []byte) {
	for i := range dek {
		dek[i] = 0
	}
}

// DeriveRow derives a unique DEK from the KEK using HKDF-SHA256 for a specific entity row.
// The derived DEK is deterministic: same (KEK, entityType, entityID) always produces the same DEK.
// Nothing needs to be stored in the database — the DEK is re-derived on every access.
// Always defer Close() to zero the DEK from memory.
func (s *AES) DeriveRow(entityType string, entityID int64, aad []byte) (gcrypt.RowEncryptorInterface, error) {
	info := []byte(fmt.Sprintf("%s:%d", entityType, entityID))
	reader := hkdf.New(sha256.New, s.kek, nil, info)
	dek := make([]byte, s.dekByteSize)
	if _, err := io.ReadFull(reader, dek); err != nil {
		zeroDEK(dek)
		return nil, err
	}
	return newRowEncryptor(dek, aad)
}

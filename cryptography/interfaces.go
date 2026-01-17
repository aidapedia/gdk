package cryptography

import (
	gencryption "github.com/aidapedia/gdk/cryptography/encryption"
)

type HashInterface interface {
	Hash(s string) string
	Verify(s string, hashed string) bool
}

type EncryptionInterface interface {
	Encrypt(plaintextDEK []byte) ([]byte, error)
	Decrypt(wrappedDEK []byte) ([]byte, error)
	// EncryptRecord encrypts the PII using the provided DEK.
	EncryptRecord(pii []byte) (*gencryption.EncryptedRecord, error)
	DecryptRecord(record *gencryption.EncryptedRecord) (string, error)
}

type TokenInterface interface {
	Encrypt(body map[string]interface{}) (string, error)
	Decrypt(token string) (map[string]interface{}, error)
}

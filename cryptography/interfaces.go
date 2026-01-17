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
	// EncryptRecord encrypts the PII using a freshly generated DEK.
	// aad is Additional Authenticated Data bound to the ciphertext (e.g. "users:123:email")
	// to prevent ciphertext swapping between rows. Pass nil to skip.
	EncryptRecord(pii []byte, aad []byte) (*gencryption.EncryptedRecord, error)
	DecryptRecord(record *gencryption.EncryptedRecord, aad []byte) (string, error)
	// DeriveRow derives a unique DEK from the KEK using HKDF for a specific entity row.
	// The DEK is deterministic (same inputs → same DEK) so nothing needs to be stored.
	// entityType is e.g. "users", entityID is the row's primary key.
	// aad is bound to GCM ciphertext to prevent cross-row swapping (typically "entityType:entityID").
	// Always defer Close() to zero the DEK from memory.
	DeriveRow(entityType string, entityID int64, aad []byte) (RowEncryptorInterface, error)
}

// RowEncryptorInterface encrypts/decrypts multiple fields with a single derived DEK.
type RowEncryptorInterface interface {
	// EncryptField encrypts a single plaintext field value and returns base64 ciphertext.
	EncryptField(plaintext []byte) (string, error)
	// DecryptField decrypts a base64 ciphertext field value and returns the plaintext.
	DecryptField(ciphertextBase64 string) (string, error)
	// Close zeros out the plaintext DEK from memory. Always defer this.
	Close()
}

type TokenInterface interface {
	Encrypt(body map[string]interface{}) (string, error)
	Decrypt(token string) (map[string]interface{}, error)
}

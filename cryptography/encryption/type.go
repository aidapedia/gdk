package encryption

// EncryptedRecord holds the components you must save to your database.
//
// For Approach A (inline), map these fields to your table columns:
//
//	email_enc   TEXT NOT NULL       → Ciphertext
//	wrapped_dek TEXT NOT NULL       → WrappedDEK
//	kek_version INT  NOT NULL       → KEKVersion
type EncryptedRecord struct {
	Ciphertext string // The encrypted PII (base64-encoded nonce+ciphertext)
	WrappedDEK string // The DEK, encrypted by your Master Key (KEK) (base64-encoded)
	KEKVersion int    // Version of the KEK used to wrap the DEK (for key rotation)
	Algorithm  string // Encryption algorithm identifier, e.g. "AES-256-GCM"
}

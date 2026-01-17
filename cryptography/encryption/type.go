package encryption

// EncryptedRecord holds the components you must save to your database.
type EncryptedRecord struct {
	Ciphertext string // The encrypted PII
	WrappedDEK string // The DEK, encrypted by your Master Key (KEK)
}

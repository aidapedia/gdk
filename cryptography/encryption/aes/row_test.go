package aes

import (
	"crypto/rand"
	"io"
	"testing"
)

func TestRowEncryptor_EncryptDecryptMultipleFields(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	enc, _ := NewAES(kek, 32, 1)

	aad := []byte("users:1")

	// Encrypt multiple fields with one derived DEK
	row, err := enc.DeriveRow("users", 1, aad)
	if err != nil {
		t.Fatalf("DeriveRow failed: %v", err)
	}

	emailEnc, err := row.EncryptField([]byte("john@example.com"))
	if err != nil {
		t.Fatalf("EncryptField email failed: %v", err)
	}
	phoneEnc, err := row.EncryptField([]byte("+6281234567890"))
	if err != nil {
		t.Fatalf("EncryptField phone failed: %v", err)
	}
	nameEnc, err := row.EncryptField([]byte("John Doe"))
	if err != nil {
		t.Fatalf("EncryptField name failed: %v", err)
	}
	row.Close()

	// Decrypt by re-deriving the same DEK (no stored wrappedDEK needed)
	row2, err := enc.DeriveRow("users", 1, aad)
	if err != nil {
		t.Fatalf("DeriveRow (decrypt) failed: %v", err)
	}
	defer row2.Close()

	email, err := row2.DecryptField(emailEnc)
	if err != nil {
		t.Fatalf("DecryptField email failed: %v", err)
	}
	phone, err := row2.DecryptField(phoneEnc)
	if err != nil {
		t.Fatalf("DecryptField phone failed: %v", err)
	}
	name, err := row2.DecryptField(nameEnc)
	if err != nil {
		t.Fatalf("DecryptField name failed: %v", err)
	}

	if email != "john@example.com" {
		t.Errorf("Expected john@example.com, got %s", email)
	}
	if phone != "+6281234567890" {
		t.Errorf("Expected +6281234567890, got %s", phone)
	}
	if name != "John Doe" {
		t.Errorf("Expected John Doe, got %s", name)
	}
}

func TestRowEncryptor_WrongAADFails(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	enc, _ := NewAES(kek, 32, 1)

	row, _ := enc.DeriveRow("users", 1, []byte("users:1"))
	ct, _ := row.EncryptField([]byte("secret"))
	row.Close()

	// Derive same DEK but open with wrong AAD
	row2, _ := enc.DeriveRow("users", 1, []byte("users:999"))
	defer row2.Close()

	_, err := row2.DecryptField(ct)
	if err == nil {
		t.Error("Expected error when decrypting with wrong AAD, got nil")
	}
}

func TestRowEncryptor_DifferentEntityIDsDifferentDEK(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	enc, _ := NewAES(kek, 32, 1)

	// Encrypt with entity ID 1
	row1, _ := enc.DeriveRow("users", 1, nil)
	ct1, _ := row1.EncryptField([]byte("same value"))
	row1.Close()

	// Try to decrypt with entity ID 2 — different derived DEK, should fail
	row2, _ := enc.DeriveRow("users", 2, nil)
	defer row2.Close()

	_, err := row2.DecryptField(ct1)
	if err == nil {
		t.Error("Expected error when decrypting with different entity ID, got nil")
	}
}

func TestRowEncryptor_DeterministicDerivation(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	enc, _ := NewAES(kek, 32, 1)

	// Derive twice with same params — should produce same DEK
	row1, _ := enc.DeriveRow("users", 42, nil)
	re1 := row1.(*RowEncryptor)
	dek1 := make([]byte, len(re1.dek))
	copy(dek1, re1.dek)
	row1.Close()

	row2, _ := enc.DeriveRow("users", 42, nil)
	re2 := row2.(*RowEncryptor)
	dek2 := make([]byte, len(re2.dek))
	copy(dek2, re2.dek)
	row2.Close()

	for i := range dek1 {
		if dek1[i] != dek2[i] {
			t.Fatal("Same (KEK, entityType, entityID) should derive the same DEK")
		}
	}
}

func TestRowEncryptor_SameFieldDifferentCiphertext(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	enc, _ := NewAES(kek, 32, 1)

	row, _ := enc.DeriveRow("users", 1, nil)
	defer row.Close()

	ct1, _ := row.EncryptField([]byte("same value"))
	ct2, _ := row.EncryptField([]byte("same value"))

	if ct1 == ct2 {
		t.Error("Same plaintext should produce different ciphertexts (unique nonce per call)")
	}
}

func TestRowEncryptor_CloseZeroesDEK(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	enc, _ := NewAES(kek, 32, 1)

	row, _ := enc.DeriveRow("users", 1, nil)
	re := row.(*RowEncryptor)
	row.Close()

	allZero := true
	for _, b := range re.dek {
		if b != 0 {
			allZero = false
			break
		}
	}
	if !allZero {
		t.Error("DEK was not zeroed after Close()")
	}
}

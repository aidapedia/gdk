package aes

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"testing"
)

func TestAES_EncryptDecrypt(t *testing.T) {
	// Generate a random KEK (32 bytes for AES-256)
	kek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, kek); err != nil {
		t.Fatal(err)
	}

	AES := NewAES(kek, 32)

	originalData := []byte("hello world")

	// Test Encrypt
	record, err := AES.EncryptRecord(originalData)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if record.Ciphertext == "" {
		t.Error("Ciphertext is empty")
	}
	if record.WrappedDEK == "" {
		t.Error("WrappedDEK is empty")
	}

	// Test Decrypt
	decryptedData, err := AES.DecryptRecord(record)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decryptedData != string(originalData) {
		t.Errorf("Expected %s, got %s", string(originalData), decryptedData)
	}
}

func TestAES_EncryptDecryptDEK(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	AES := NewAES(kek, 32)

	dek := make([]byte, 32)
	io.ReadFull(rand.Reader, dek)

	wrapped, err := AES.Encrypt(dek)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	unwrapped, err := AES.Decrypt(wrapped)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(dek, unwrapped) {
		t.Error("Unwrapped DEK does not match original DEK")
	}
}

func TestAES_InvalidKEK(t *testing.T) {
	// AES keys must be 16, 24, or 32 bytes.
	invalidKEK := []byte("too short")
	AES := NewAES(invalidKEK, 32)

	dek := make([]byte, 32)
	io.ReadFull(rand.Reader, dek)

	_, err := AES.Encrypt(dek)
	if err == nil {
		t.Error("Expected error for invalid KEK size, got nil")
	}
}

func TestAES_CorruptedCiphertext(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	AES := NewAES(kek, 32)

	originalData := []byte("secret")
	record, _ := AES.EncryptRecord(originalData)

	// Corrupt the ciphertext
	decodedData, _ := base64.StdEncoding.DecodeString(record.Ciphertext)
	decodedData[len(decodedData)-1] ^= 0xFF
	record.Ciphertext = base64.StdEncoding.EncodeToString(decodedData)

	_, err := AES.DecryptRecord(record)
	if err == nil {
		t.Error("Expected error for corrupted ciphertext, got nil")
	}
}

func BenchmarkAES_Encrypt(b *testing.B) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	AES := NewAES(kek, 32)
	pii := []byte("highly sensitive information")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AES.Encrypt(pii)
	}
}

func BenchmarkAES_Decrypt(b *testing.B) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	AES := NewAES(kek, 32)
	pii := []byte("highly sensitive information")
	record, _ := AES.Encrypt(pii)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AES.Decrypt(record)
	}
}

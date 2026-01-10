package aes

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"testing"
)

func TestShield_EnvelopeEncryptDecrypt(t *testing.T) {
	// Generate a random KEK (32 bytes for AES-256)
	kek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, kek); err != nil {
		t.Fatal(err)
	}

	shield := NewShield(kek)

	originalData := []byte("hello world")

	// Test EnvelopeEncrypt
	record, err := shield.EnvelopeEncrypt(originalData)
	if err != nil {
		t.Fatalf("EnvelopeEncrypt failed: %v", err)
	}

	if record.Ciphertext == "" {
		t.Error("Ciphertext is empty")
	}
	if record.WrappedDEK == "" {
		t.Error("WrappedDEK is empty")
	}

	// Test EnvelopeDecrypt
	decryptedData, err := shield.EnvelopeDecrypt(record)
	if err != nil {
		t.Fatalf("EnvelopeDecrypt failed: %v", err)
	}

	if decryptedData != string(originalData) {
		t.Errorf("Expected %s, got %s", string(originalData), decryptedData)
	}
}

func TestShield_EncryptDecryptDEK(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	shield := NewShield(kek)

	dek := make([]byte, 32)
	io.ReadFull(rand.Reader, dek)

	wrapped, err := shield.encrypt(dek)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	unwrapped, err := shield.decrypt(wrapped)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(dek, unwrapped) {
		t.Error("Unwrapped DEK does not match original DEK")
	}
}

func TestShield_InvalidKEK(t *testing.T) {
	// AES keys must be 16, 24, or 32 bytes.
	invalidKEK := []byte("too short")
	shield := NewShield(invalidKEK)

	dek := make([]byte, 32)
	io.ReadFull(rand.Reader, dek)

	_, err := shield.encrypt(dek)
	if err == nil {
		t.Error("Expected error for invalid KEK size, got nil")
	}
}

func TestShield_CorruptedCiphertext(t *testing.T) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	shield := NewShield(kek)

	originalData := []byte("secret")
	record, _ := shield.EnvelopeEncrypt(originalData)

	// Corrupt the ciphertext
	decodedData, _ := base64.StdEncoding.DecodeString(record.Ciphertext)
	decodedData[len(decodedData)-1] ^= 0xFF
	record.Ciphertext = base64.StdEncoding.EncodeToString(decodedData)

	_, err := shield.EnvelopeDecrypt(record)
	if err == nil {
		t.Error("Expected error for corrupted ciphertext, got nil")
	}
}

func BenchmarkShield_EnvelopeEncrypt(b *testing.B) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	shield := NewShield(kek)
	pii := []byte("highly sensitive information")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = shield.EnvelopeEncrypt(pii)
	}
}

func BenchmarkShield_EnvelopeDecrypt(b *testing.B) {
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	shield := NewShield(kek)
	pii := []byte("highly sensitive information")
	record, _ := shield.EnvelopeEncrypt(pii)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = shield.EnvelopeDecrypt(record)
	}
}

package aes

import (
	"testing"
)

func TestShield_EncryptDecrypt(t *testing.T) {
	passphrase := "super-secret-password"
	shield := NewShield(passphrase)
	plainText := "Hello, World! This is a secret message."

	// Test Encrypt
	encrypted, err := shield.Encrypt(plainText)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted == plainText {
		t.Error("Encrypted text should not be equal to plain text")
	}

	// Test Decrypt
	decrypted, err := shield.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plainText {
		t.Errorf("Decrypted text = %v, want %v", decrypted, plainText)
	}
}

func TestShield_DecryptWithWrongKey(t *testing.T) {
	passphrase1 := "password-1"
	passphrase2 := "password-2"
	plainText := "Secret message"

	shield1 := NewShield(passphrase1)
	shield2 := NewShield(passphrase2)

	encrypted, err := shield1.Encrypt(plainText)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Try to decrypt with wrong shield (different key)
	_, err = shield2.Decrypt(encrypted)
	if err == nil {
		t.Error("Expected error when decrypting with wrong key, got nil")
	}
}

func TestShield_DecryptTampered(t *testing.T) {
	passphrase := "password"
	shield := NewShield(passphrase)
	plainText := "Top Secret"

	encrypted, err := shield.Encrypt(plainText)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Tamper with the encrypted string (Base64)
	bytes := []byte(encrypted)
	if len(bytes) > 10 {
		// Change a character in the middle
		if bytes[10] == 'A' {
			bytes[10] = 'B'
		} else {
			bytes[10] = 'A'
		}
	}
	tampered := string(bytes)

	_, err = shield.Decrypt(tampered)
	if err == nil {
		t.Error("Expected error when decrypting tampered data, got nil")
	}
}

func TestShield_InvalidBase64(t *testing.T) {
	shield := NewShield("password")
	_, err := shield.Decrypt("invalid-base64-!@#$%^&*")
	if err == nil {
		t.Error("Expected error for invalid base64 input, got nil")
	}
}

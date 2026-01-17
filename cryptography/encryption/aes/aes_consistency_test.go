package aes

import (
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func TestWrappedDEKConsistency(t *testing.T) {
	// 1. Setup KEK
	kek := make([]byte, 32)
	io.ReadFull(rand.Reader, kek)
	AES, _ := NewAES(kek, 16, 1)

	pii := []byte("secret data")

	// 2. Encrypt twice
	record1, _ := AES.EncryptRecord(pii, nil)
	record2, _ := AES.EncryptRecord(pii, nil)

	fmt.Println("\n--- WrappedDEK Comparison ---")
	fmt.Printf("Attempt 1 WrappedDEK: %s\n", record1.WrappedDEK)
	fmt.Printf("Attempt 2 WrappedDEK: %s\n", record2.WrappedDEK)

	if record1.WrappedDEK == record2.WrappedDEK {
		fmt.Println("RESULT: They are the SAME (Unexpected!)")
	} else {
		fmt.Println("RESULT: They are DIFFERENT (This is correct/secure)")
	}
	fmt.Println("----------------------------")
}

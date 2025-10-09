package cryptography

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash hashes the string using fnv64a and base64
func Hash(s string) string {
	password := []byte(s)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

// CheckHash checks the string with the hashed string
func CheckHash(s string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(s))
	if err != nil {
		return false
	}
	return true
}

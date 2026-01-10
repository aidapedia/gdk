package bcrypt

import (
	"github.com/aidapedia/gdk/cryptography/hash"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
	cost int
}

func New(cost int) hash.Interface {
	return &Bcrypt{
		cost: cost,
	}
}

// Hash hashes the string using fnv64a and base64
func (b *Bcrypt) Hash(s string) string {
	password := []byte(s)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, b.cost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

// Verify checks the string with the hashed string
func (b *Bcrypt) Verify(s string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(s))
	if err != nil {
		return false
	}
	return true
}

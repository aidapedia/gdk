package validation

import "net/mail"

// IsEmail checks if the given string is a valid email address
func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

package validation

import "regexp"

// ValidatePasswordStrength is a function to validate password strength
func ValidatePasswordStrength(password string) (bool, error) {
	return regexp.Match(`/^(?=[^a-z]*[a-z])(?=[^A-Z]*[A-Z])(?=\D*\d)(?=[^!#%]*[!#%])[A-Za-z0-9!#%]{8,32}$/`, []byte(password))
}

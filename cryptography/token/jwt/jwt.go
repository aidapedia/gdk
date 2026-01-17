package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	gcrypt "github.com/aidapedia/gdk/cryptography"
)

// JWT implements TokenInterface
type JWT struct {
	key          Key
	parserOption []jwt.ParserOption
}

// New creates a new JWT instance
func New(key Key, parseOpt ...jwt.ParserOption) gcrypt.TokenInterface {
	return &JWT{
		key:          key,
		parserOption: parseOpt,
	}
}

// Sign signs a token using a private key
func (s *JWT) Encrypt(claims map[string]interface{}) (string, error) {
	key, err := s.key.GetEncryptKey()
	if err != nil {
		return "", fmt.Errorf("generate token parse key error: %w", err)
	}
	return jwt.NewWithClaims(s.key.GetSigningMethod(), jwt.MapClaims(claims)).
		SignedString(key)
}

// VerifyToken verifies a token using a public key
func (s *JWT) Decrypt(token string) (map[string]interface{}, error) {
	key, err := s.key.GetDecryptKey()
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		err := s.key.Validate(token)
		if err != nil {
			return nil, err
		}
		return key, nil
	}, s.parserOption...)
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate invalid")
	}

	return claims, nil
}

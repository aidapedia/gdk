package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWTToken is a struct that holds the private keys
type jwtToken struct {
	privatekey []byte
	publickey  []byte
	validator  jwt.Validator
}

var JWT jwtToken

func New(privatekey, publickey []byte, parseOpt ...jwt.ParserOption) {
	JWT = jwtToken{
		privatekey: privatekey,
		publickey:  publickey,
		validator: *jwt.NewValidator(
			parseOpt...,
		),
	}
}

// SignToken signs a token using a private key
func SignToken(body map[string]interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(JWT.privatekey)
	if err != nil {
		return "", fmt.Errorf("generate token parse key error: %w", err)
	}
	token := jwt.New(jwt.SigningMethodRS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	for k, v := range body {
		claims[k] = v
	}
	// Generate encoded token and send it as response.
	return token.SignedString(key)
}

// VerifyToken verifies a token using a public key
func VerifyToken(token string) (map[string]interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(JWT.publickey)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	err = JWT.validator.Validate(tok.Claims)
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate invalid")
	}

	return claims, nil
}

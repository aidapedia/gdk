package jwt

import "github.com/golang-jwt/jwt/v5"

type Key interface {
	Validate(token *jwt.Token) error
	GetEncryptKey() (any, error)
	GetDecryptKey() (any, error)
	GetSigningMethod() jwt.SigningMethod
}

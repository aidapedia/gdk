package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type SingleKey struct {
	key []byte

	method     jwt.SigningMethod
	methodType SingleKeyType
}

// NewSingleKey creates a new SingleKey instance.
// SigningMethod eligible: HS256, HS384, HS512
func NewSingleKey(key []byte, signingMethod jwt.SigningMethod) (Key, error) {
	methodType := determineSingleKeyType(signingMethod)
	if methodType == "" {
		return nil, fmt.Errorf("invalid signing method type")
	}

	return &SingleKey{
		key:        key,
		method:     signingMethod,
		methodType: methodType,
	}, nil
}

func (s *SingleKey) Validate(token *jwt.Token) error {
	switch s.methodType {
	case HMACKeyType:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	default:
		return fmt.Errorf("unexpected signing method: %v", s.methodType)
	}

	return nil
}

func (s *SingleKey) GetSigningMethod() jwt.SigningMethod {
	return s.method
}

func (s *SingleKey) GetEncryptKey() (any, error) {
	switch s.methodType {
	case HMACKeyType:
		return s.key, nil
	default:
		return nil, fmt.Errorf("unknown key type")
	}
}

func (s *SingleKey) GetDecryptKey() (any, error) {
	switch s.methodType {
	case HMACKeyType:
		return s.key, nil
	default:
		return nil, fmt.Errorf("unknown key type")
	}
}

type SingleKeyType string

const (
	HMACKeyType SingleKeyType = "HMAC"
)

func determineSingleKeyType(signingMethod jwt.SigningMethod) SingleKeyType {
	switch signingMethod {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		return HMACKeyType
	default:
		return ""
	}
}

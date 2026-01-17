package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type PairKey struct {
	privatekey []byte
	publickey  []byte

	method     jwt.SigningMethod
	methodType PairKeyType
}

// NewPairKey creates a new PairKey instance.
// SigningMethod eligible: RS256, RS384, RS512, PS256, PS384, PS512, ES256, ES384, ES512, EdDSA
func NewPairKey(privatekey, publickey []byte, signingMethod jwt.SigningMethod) (Key, error) {
	methodType := determinePairKeyType(signingMethod)
	if methodType == "" {
		return nil, fmt.Errorf("invalid signing method type")
	}
	return &PairKey{
		privatekey: privatekey,
		publickey:  publickey,
		method:     signingMethod,
		methodType: methodType,
	}, nil
}

func (s *PairKey) Validate(token *jwt.Token) error {
	switch s.methodType {
	case RSAKeyType:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	case RSAPSSKeyType:
		if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
			return fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	case ECDSAKeyType:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	case ED25519KeyType:
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	default:
		return fmt.Errorf("unexpected signing method: %v", s.methodType)
	}

	return nil
}

func (s *PairKey) GetSigningMethod() jwt.SigningMethod {
	return s.method
}

func (s *PairKey) GetEncryptKey() (any, error) {
	switch s.methodType {
	case RSAKeyType, RSAPSSKeyType:
		return jwt.ParseRSAPrivateKeyFromPEM(s.privatekey)
	case ECDSAKeyType:
		return jwt.ParseECPrivateKeyFromPEM(s.privatekey)
	case ED25519KeyType:
		return jwt.ParseEdPrivateKeyFromPEM(s.privatekey)
	default:
		return nil, fmt.Errorf("unknown key type")
	}
}

func (s *PairKey) GetDecryptKey() (any, error) {
	switch s.methodType {
	case RSAKeyType, RSAPSSKeyType:
		return jwt.ParseRSAPublicKeyFromPEM(s.publickey)
	case ECDSAKeyType:
		return jwt.ParseECPublicKeyFromPEM(s.publickey)
	case ED25519KeyType:
		return jwt.ParseEdPublicKeyFromPEM(s.publickey)
	default:
		return nil, fmt.Errorf("unknown key type")
	}
}

type PairKeyType string

const (
	RSAKeyType     PairKeyType = "RSA"
	RSAPSSKeyType  PairKeyType = "RSAPSS"
	ECDSAKeyType   PairKeyType = "ECDSA"
	ED25519KeyType PairKeyType = "ED25519"
)

func determinePairKeyType(signingMethod jwt.SigningMethod) PairKeyType {
	switch signingMethod {
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		return RSAKeyType
	case jwt.SigningMethodPS256, jwt.SigningMethodPS384, jwt.SigningMethodPS512:
		return RSAPSSKeyType
	case jwt.SigningMethodEdDSA:
		return ED25519KeyType
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		return ECDSAKeyType
	default:
		return ""
	}
}

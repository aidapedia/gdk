package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func generateRSAPEM(t *testing.T) ([]byte, []byte) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return privPEM, pubPEM
}

func generateECDSAPEM(t *testing.T, curve elliptic.Curve) ([]byte, []byte) {
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ecdsa key: %v", err)
	}

	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal ecdsa key: %v", err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return privPEM, pubPEM
}

func generateEd25519PEM(t *testing.T) ([]byte, []byte) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ed25519 key: %v", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatalf("failed to marshal ed25519 key: %v", err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return privPEM, pubPEM
}

func TestJWT_SigningMethods(t *testing.T) {
	claims := map[string]interface{}{
		"user_id": "12345",
		"role":    "admin",
	}

	rsaPriv, rsaPub := generateRSAPEM(t)
	ecdsa256Priv, ecdsa256Pub := generateECDSAPEM(t, elliptic.P256())
	ecdsa384Priv, ecdsa384Pub := generateECDSAPEM(t, elliptic.P384())
	ecdsa521Priv, ecdsa521Pub := generateECDSAPEM(t, elliptic.P521())
	ed25519Priv, ed25519Pub := generateEd25519PEM(t)

	tests := []struct {
		name          string
		signingMethod jwt.SigningMethod
		isPair        bool
		privKey       []byte
		pubKey        []byte
		hmacKey       []byte
		hmacMsg       []byte
	}{
		// Single Key (HMAC)
		{
			name:          "HS256",
			signingMethod: jwt.SigningMethodHS256,
			isPair:        false,
			hmacKey:       []byte("secret"),
			hmacMsg:       []byte("message"),
		},
		{
			name:          "HS384",
			signingMethod: jwt.SigningMethodHS384,
			isPair:        false,
			hmacKey:       []byte("secret"),
			hmacMsg:       []byte("message"),
		},
		{
			name:          "HS512",
			signingMethod: jwt.SigningMethodHS512,
			isPair:        false,
			hmacKey:       []byte("secret"),
			hmacMsg:       []byte("message"),
		},
		// Pair Key (RSA)
		{
			name:          "RS256",
			signingMethod: jwt.SigningMethodRS256,
			isPair:        true,
			privKey:       rsaPriv,
			pubKey:        rsaPub,
		},
		{
			name:          "RS384",
			signingMethod: jwt.SigningMethodRS384,
			isPair:        true,
			privKey:       rsaPriv,
			pubKey:        rsaPub,
		},
		{
			name:          "RS512",
			signingMethod: jwt.SigningMethodRS512,
			isPair:        true,
			privKey:       rsaPriv,
			pubKey:        rsaPub,
		},
		// Pair Key (RSAPSS)
		{
			name:          "PS256",
			signingMethod: jwt.SigningMethodPS256,
			isPair:        true,
			privKey:       rsaPriv,
			pubKey:        rsaPub,
		},
		{
			name:          "PS384",
			signingMethod: jwt.SigningMethodPS384,
			isPair:        true,
			privKey:       rsaPriv,
			pubKey:        rsaPub,
		},
		{
			name:          "PS512",
			signingMethod: jwt.SigningMethodPS512,
			isPair:        true,
			privKey:       rsaPriv,
			pubKey:        rsaPub,
		},
		// Pair Key (ECDSA)
		{
			name:          "ES256",
			signingMethod: jwt.SigningMethodES256,
			isPair:        true,
			privKey:       ecdsa256Priv,
			pubKey:        ecdsa256Pub,
		},
		{
			name:          "ES384",
			signingMethod: jwt.SigningMethodES384,
			isPair:        true,
			privKey:       ecdsa384Priv,
			pubKey:        ecdsa384Pub,
		},
		{
			name:          "ES512",
			signingMethod: jwt.SigningMethodES512,
			isPair:        true,
			privKey:       ecdsa521Priv,
			pubKey:        ecdsa521Pub,
		},
		// Pair Key (EdDSA)
		{
			name:          "EdDSA",
			signingMethod: jwt.SigningMethodEdDSA,
			isPair:        true,
			privKey:       ed25519Priv,
			pubKey:        ed25519Pub,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key Key
			var err error

			if tt.isPair {
				key, err = NewPairKey(tt.privKey, tt.pubKey, tt.signingMethod)
			} else {
				key, err = NewSingleKey(tt.hmacKey, tt.signingMethod)
			}
			if err != nil {
				t.Fatalf("failed to create key: %v", err)
			}

			j := New(key)

			// Test Encrypt (Sign)
			token, err := j.Encrypt(claims)
			if err != nil {
				t.Fatalf("failed to encrypt: %v", err)
			}
			if token == "" {
				t.Fatal("token is empty")
			}

			// Test Decrypt (Verify)
			decryptedClaims, err := j.Decrypt(token)
			if err != nil {
				t.Fatalf("failed to decrypt: %v", err)
			}

			// Verify claims
			for k, v := range claims {
				if !reflect.DeepEqual(v, decryptedClaims[k]) {
					t.Errorf("claim %s: expected %v, got %v", k, v, decryptedClaims[k])
				}
			}
		})
	}
}

func TestNewSingleKey_InvalidMethod(t *testing.T) {
	_, err := NewSingleKey([]byte("key"), jwt.SigningMethodRS256)
	if err == nil {
		t.Error("expected error for invalid method, got nil")
	}
}

func TestNewPairKey_InvalidMethod(t *testing.T) {
	_, err := NewPairKey([]byte("priv"), []byte("pub"), jwt.SigningMethodHS256)
	if err == nil {
		t.Error("expected error for invalid method, got nil")
	}
}

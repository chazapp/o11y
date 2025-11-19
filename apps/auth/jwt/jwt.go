package jwt

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-jose/go-jose/v4"
)

// Claims represents the JWT claims used in the authentication system
type Claims struct {
	Email     string `json:"email"`
	Issuer    string `json:"iss"`
	ExpiresAt int64  `json:"exp"`
}

func LoadJWTKeys(jwtPrivateKeyPath, jwtPublicKeyPath string) (jose.JSONWebKey, jose.JSONWebKey, error) {
	var jwkPriv, jwkPub jose.JSONWebKey

	bytesPriv, err := os.ReadFile(jwtPrivateKeyPath)
	if err != nil {
		return jwkPriv, jwkPub, fmt.Errorf("failed to read JWT private key file %s: %w", jwtPrivateKeyPath, err)
	}
	bytesPub, err := os.ReadFile(jwtPublicKeyPath)
	if err != nil {
		return jwkPriv, jwkPub, fmt.Errorf("failed to read JWT public key file %s: %w", jwtPublicKeyPath, err)
	}

	err = jwkPriv.UnmarshalJSON(bytesPriv)
	if err != nil {
		return jwkPriv, jwkPub, fmt.Errorf("failed to unmarshal JWT private key: %w", err)
	}
	err = jwkPub.UnmarshalJSON(bytesPub)
	if err != nil {
		return jwkPriv, jwkPub, fmt.Errorf("failed to unmarshal JWT public key: %w", err)
	}
	return jwkPriv, jwkPub, nil
}

func GenerateJWT(duration, email string, jwkPrivate jose.JSONWebKey) (string, error) {
	d, err := time.ParseDuration(duration)

	if err != nil {
		return "", fmt.Errorf("invalid duration format: %w", err)
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: jwkPrivate.Key}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create signer: %w", err)
	}

	claims := Claims{
		Email:     email,
		Issuer:    "o11y",
		ExpiresAt: time.Now().Add(d).Unix(),
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}
	payload, err := signer.Sign(claimsJSON)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	token, err := payload.CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize JWT: %w", err)
	}
	return token, nil
}

// VerifyJWT validates a JWT token and returns the claims if valid
func VerifyJWT(tokenString string, jwkPublic jose.JSONWebKey) (*Claims, error) {
	sig, err := jose.ParseSigned(tokenString, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	payload, err := sig.Verify(jwkPublic.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token has expired")
	}
	return &claims, nil
}

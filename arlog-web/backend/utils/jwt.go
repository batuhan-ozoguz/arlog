package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// OktaJWKS represents the Okta JSON Web Key Set
type OktaJWKS struct {
	Keys []OktaJWK `json:"keys"`
}

// OktaJWK represents a single JSON Web Key
type OktaJWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// OktaClaims represents the claims in an Okta JWT token
type OktaClaims struct {
	Sub    string   `json:"sub"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
	jwt.RegisteredClaims
}

// JWTValidator handles JWT token validation
type JWTValidator struct {
	oktaDomain string
	clientID   string
	jwksCache  map[string]*rsa.PublicKey
	cacheTime  time.Time
}

// NewJWTValidator creates a new JWT validator
func NewJWTValidator(oktaDomain, clientID string) *JWTValidator {
	return &JWTValidator{
		oktaDomain: oktaDomain,
		clientID:   clientID,
		jwksCache:  make(map[string]*rsa.PublicKey),
	}
}

// ValidateToken validates an Okta JWT token
func (v *JWTValidator) ValidateToken(tokenString string) (*OktaClaims, error) {
	// Parse token without validation first to get the kid (key ID)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &OktaClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Get the key ID from the token header
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("token missing kid header")
	}

	// Get the public key for this kid
	publicKey, err := v.getPublicKey(kid)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Parse and validate the token with the public key
	claims := &OktaClaims{}
	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Validate audience (client ID)
	if !claims.VerifyAudience(v.clientID, true) {
		return nil, fmt.Errorf("invalid audience")
	}

	// Validate issuer
	expectedIssuer := fmt.Sprintf("https://%s/oauth2/default", v.oktaDomain)
	if !claims.VerifyIssuer(expectedIssuer, true) {
		return nil, fmt.Errorf("invalid issuer")
	}

	return claims, nil
}

// getPublicKey retrieves the public key for a given key ID
func (v *JWTValidator) getPublicKey(kid string) (*rsa.PublicKey, error) {
	// Check cache first (refresh every hour)
	if time.Since(v.cacheTime) < time.Hour {
		if key, ok := v.jwksCache[kid]; ok {
			return key, nil
		}
	}

	// Fetch JWKS from Okta
	jwksURL := fmt.Sprintf("https://%s/oauth2/default/v1/keys", v.oktaDomain)
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks OktaJWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	// Find the key with matching kid
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			publicKey, err := v.parsePublicKey(key)
			if err != nil {
				return nil, err
			}

			// Update cache
			v.jwksCache[kid] = publicKey
			v.cacheTime = time.Now()

			return publicKey, nil
		}
	}

	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

// parsePublicKey converts a JWK to an RSA public key
func (v *JWTValidator) parsePublicKey(jwk OktaJWK) (*rsa.PublicKey, error) {
	// Decode the modulus (n)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	// Decode the exponent (e)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert bytes to big integers
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// Create the RSA public key
	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}

// ExtractBearerToken extracts the bearer token from the Authorization header
func ExtractBearerToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return parts[1], nil
}


package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/configs"
)

var (
	// ErrInvalidToken is returned when the token format, algorithm, or signature is invalid.
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when the token expiration time has passed.
	ErrExpiredToken = errors.New("token has expired")
	// ErrNotInitialized is returned when the global JWT manager is accessed without initialization.
	ErrNotInitialized = errors.New("jwt manager has not been initialized")
)

// Claims represents the payload contained in the JWT.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
	Iss    string `json:"iss,omitempty"`
}

// TokenManager provides methods to sign and verify JWT tokens using HMAC-SHA256.
type TokenManager struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

var (
	defaultManager *TokenManager
	mu             sync.RWMutex
)

// NewTokenManager creates a new TokenManager instance configured with JWTConfig.
func NewTokenManager(cfg configs.JWTConfig) *TokenManager {
	exp := cfg.ExpirationDuration
	if exp == 0 {
		exp = 24 * time.Hour
	}

	return &TokenManager{
		secret:     []byte(cfg.Secret),
		expiration: exp,
		issuer:     cfg.Issuer,
	}
}

// Init initializes the default global token manager.
func Init(cfg configs.JWTConfig) *TokenManager {
	mu.Lock()
	defer mu.Unlock()
	defaultManager = NewTokenManager(cfg)
	return defaultManager
}

// GenerateToken generates a signed JWT token using the global manager.
func GenerateToken(userID, email, role string) (string, error) {
	mu.RLock()
	mgr := defaultManager
	mu.RUnlock()

	if mgr == nil {
		return "", ErrNotInitialized
	}
	return mgr.GenerateToken(userID, email, role)
}

// ParseToken parses and verifies a JWT token using the global manager.
func ParseToken(token string) (Claims, error) {
	mu.RLock()
	mgr := defaultManager
	mu.RUnlock()

	if mgr == nil {
		return Claims{}, ErrNotInitialized
	}
	return mgr.ParseToken(token)
}

// GenerateToken creates a signed JWT string for the given user claims.
func (m *TokenManager) GenerateToken(userID, email, role string) (string, error) {
	if len(m.secret) == 0 {
		return "", errors.New("jwt secret key cannot be empty")
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Exp:    time.Now().Add(m.expiration).Unix(),
		Iss:    m.issuer,
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	unsignedToken := encodeSegment(headerJSON) + "." + encodeSegment(claimsJSON)
	signature := m.sign(unsignedToken)

	return unsignedToken + "." + signature, nil
}

// ParseToken parses, verifies the algorithm, checks the signature, and validates claims.
func (m *TokenManager) ParseToken(token string) (Claims, error) {
	if len(m.secret) == 0 {
		return Claims{}, errors.New("jwt secret key cannot be empty")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	// 1. Decode and verify header algorithm explicitly to prevent 'alg: none' and algorithm confusion
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Claims{}, ErrInvalidToken
	}
	if !strings.EqualFold(header.Alg, "HS256") {
		return Claims{}, fmt.Errorf("%w: unexpected signing algorithm '%s', expected 'HS256'", ErrInvalidToken, header.Alg)
	}

	// 2. Verify signature
	unsignedToken := parts[0] + "." + parts[1]
	expectedSig := m.sign(unsignedToken)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return Claims{}, ErrInvalidToken
	}

	// 3. Decode and parse claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}

	// 4. Validate expiration
	if claims.Exp < time.Now().Unix() {
		return Claims{}, ErrExpiredToken
	}

	// 5. Validate UserID claim
	if strings.TrimSpace(claims.UserID) == "" {
		return Claims{}, fmt.Errorf("%w: user_id claim cannot be empty", ErrInvalidToken)
	}

	// 6. Validate issuer if configured on manager and present in token
	if m.issuer != "" && claims.Iss != "" && claims.Iss != m.issuer {
		return Claims{}, fmt.Errorf("%w: token issuer mismatch", ErrInvalidToken)
	}

	return claims, nil
}

func (m *TokenManager) sign(data string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(data))
	return encodeSegment(mac.Sum(nil))
}

func encodeSegment(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

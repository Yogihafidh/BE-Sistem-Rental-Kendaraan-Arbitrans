package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/configs"
)

func TestTokenManager_GenerateAndParse(t *testing.T) {
	cfg := configs.JWTConfig{
		Secret:             "my-very-long-and-secure-test-secret-key",
		ExpirationDuration: 1 * time.Hour,
		Issuer:             "arbitrans-test",
	}

	mgr := NewTokenManager(cfg)

	token, err := mgr.GenerateToken("user-123", "user@example.com", "admin")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	claims, err := mgr.ParseToken(token)
	if err != nil {
		t.Fatalf("unexpected error parsing token: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("expected UserID 'user-123', got '%s'", claims.UserID)
	}
	if claims.Email != "user@example.com" {
		t.Errorf("expected Email 'user@example.com', got '%s'", claims.Email)
	}
	if claims.Role != "admin" {
		t.Errorf("expected Role 'admin', got '%s'", claims.Role)
	}
	if claims.Iss != "arbitrans-test" {
		t.Errorf("expected Iss 'arbitrans-test', got '%s'", claims.Iss)
	}
}

func TestTokenManager_ExpiredToken(t *testing.T) {
	cfg := configs.JWTConfig{
		Secret:             "my-very-long-and-secure-test-secret-key",
		ExpirationDuration: -1 * time.Hour, // already expired
		Issuer:             "arbitrans-test",
	}

	mgr := NewTokenManager(cfg)

	token, err := mgr.GenerateToken("user-123", "user@example.com", "admin")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = mgr.ParseToken(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("expected ErrExpiredToken, got: %v", err)
	}
}

func TestTokenManager_InvalidSignature(t *testing.T) {
	mgr1 := NewTokenManager(configs.JWTConfig{Secret: "secret-key-one-12345678901234567890", ExpirationDuration: time.Hour})
	mgr2 := NewTokenManager(configs.JWTConfig{Secret: "secret-key-two-12345678901234567890", ExpirationDuration: time.Hour})

	token, err := mgr1.GenerateToken("user-123", "user@example.com", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = mgr2.ParseToken(token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got: %v", err)
	}
}

func TestTokenManager_RejectAlgNone(t *testing.T) {
	mgr := NewTokenManager(configs.JWTConfig{Secret: "secret-key-one-12345678901234567890", ExpirationDuration: time.Hour})

	// Manually construct token with alg: none
	header := map[string]string{"alg": "none", "typ": "JWT"}
	claims := Claims{UserID: "attacker", Email: "hacker@evil.com", Role: "admin", Exp: time.Now().Add(time.Hour).Unix()}

	hBytes, _ := json.Marshal(header)
	cBytes, _ := json.Marshal(claims)

	fakeToken := base64.RawURLEncoding.EncodeToString(hBytes) + "." + base64.RawURLEncoding.EncodeToString(cBytes) + "."

	_, err := mgr.ParseToken(fakeToken)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for alg: none, got: %v", err)
	}
}

func TestTokenManager_RejectUnexpectedAlgorithm(t *testing.T) {
	mgr := NewTokenManager(configs.JWTConfig{Secret: "secret-key-one-12345678901234567890", ExpirationDuration: time.Hour})

	header := map[string]string{"alg": "RS256", "typ": "JWT"}
	claims := Claims{UserID: "attacker", Email: "hacker@evil.com", Role: "admin", Exp: time.Now().Add(time.Hour).Unix()}

	hBytes, _ := json.Marshal(header)
	cBytes, _ := json.Marshal(claims)

	fakeToken := base64.RawURLEncoding.EncodeToString(hBytes) + "." + base64.RawURLEncoding.EncodeToString(cBytes) + ".fakesig"

	_, err := mgr.ParseToken(fakeToken)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for alg: RS256, got: %v", err)
	}
}

func TestTokenManager_RejectMissingUserID(t *testing.T) {
	mgr := NewTokenManager(configs.JWTConfig{Secret: "secret-key-one-12345678901234567890", ExpirationDuration: time.Hour})

	token, err := mgr.GenerateToken("", "user@example.com", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = mgr.ParseToken(token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for empty user_id, got: %v", err)
	}
}

func TestTokenManager_RejectIssuerMismatch(t *testing.T) {
	mgr := NewTokenManager(configs.JWTConfig{Secret: "secret-key-one-12345678901234567890", ExpirationDuration: time.Hour, Issuer: "expected-issuer"})

	// Token generated with different issuer
	otherMgr := NewTokenManager(configs.JWTConfig{Secret: "secret-key-one-12345678901234567890", ExpirationDuration: time.Hour, Issuer: "different-issuer"})
	token, _ := otherMgr.GenerateToken("user-1", "user@test.com", "staff")

	_, err := mgr.ParseToken(token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for issuer mismatch, got: %v", err)
	}
}

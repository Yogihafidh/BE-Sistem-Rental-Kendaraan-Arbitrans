package configs

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear environment variables
	os.Clearenv()
	_ = os.Setenv("JWT_SECRET", "my-secure-development-secret-key-123")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.App.Env != "development" {
		t.Errorf("expected App.Env 'development', got '%s'", cfg.App.Env)
	}
	if cfg.App.Port != "8080" {
		t.Errorf("expected App.Port '8080', got '%s'", cfg.App.Port)
	}
	if cfg.DB.MaxConns != 25 {
		t.Errorf("expected DB.MaxConns 25, got %d", cfg.DB.MaxConns)
	}
	if cfg.JWT.ExpirationDuration != 24*time.Hour {
		t.Errorf("expected JWT.ExpirationDuration 24h, got %v", cfg.JWT.ExpirationDuration)
	}
}

func TestValidate_MissingJWTSecret(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "development", Port: "8080"},
		DB:  DatabaseConfig{URL: "postgres://localhost/test", MaxConns: 10, MinConns: 2},
		JWT: JWTConfig{Secret: ""},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for empty JWT_SECRET, got nil")
	}
}

func TestValidate_InvalidAppEnv(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "unknown_environment", Port: "8080"},
		DB:  DatabaseConfig{URL: "postgres://localhost/test", MaxConns: 10, MinConns: 2},
		JWT: JWTConfig{Secret: "test-secret", ExpirationDuration: time.Hour},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid APP_ENV, got nil")
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "development", Port: "999999"},
		DB:  DatabaseConfig{URL: "postgres://localhost/test", MaxConns: 10, MinConns: 2},
		JWT: JWTConfig{Secret: "test-secret", ExpirationDuration: time.Hour},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid APP_PORT, got nil")
	}
}

func TestValidate_InvalidDBPort(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "development", Port: "8080"},
		DB: DatabaseConfig{
			Host:     "localhost",
			Port:     99999, // invalid port
			User:     "postgres",
			Name:     "rental",
			MaxConns: 10,
			MinConns: 2,
		},
		JWT: JWTConfig{Secret: "test-secret", ExpirationDuration: time.Hour},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid DB_PORT, got nil")
	}
}

func TestValidate_ProductionInsecureSecret(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "production", Port: "8080"},
		DB:  DatabaseConfig{URL: "postgres://localhost/test", MaxConns: 10, MinConns: 2},
		JWT: JWTConfig{Secret: "supersecretkey", ExpirationDuration: time.Hour},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for weak JWT_SECRET in production, got nil")
	}
}

func TestValidate_ProductionCORSWildcardWithCredentials(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "production", Port: "8080"},
		DB:  DatabaseConfig{URL: "postgres://localhost/test", MaxConns: 10, MinConns: 2},
		JWT: JWTConfig{Secret: "a-secure-production-secret-with-32-chars!", ExpirationDuration: time.Hour},
		CORS: CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for wildcard CORS origin with AllowCredentials in production, got nil")
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	dbWithURL := DatabaseConfig{URL: "postgres://user:pass@remote:5432/mydb?sslmode=require"}
	if dbWithURL.GetDSN() != "postgres://user:pass@remote:5432/mydb?sslmode=require" {
		t.Errorf("expected URL DSN, got: %s", dbWithURL.GetDSN())
	}

	dbWithParams := DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "myuser",
		Password: "mypassword",
		Name:     "rental",
		SSLMode:  "disable",
	}
	expected := "postgres://myuser:mypassword@127.0.0.1:5432/rental?sslmode=disable"
	if dbWithParams.GetDSN() != expected {
		t.Errorf("expected '%s', got '%s'", expected, dbWithParams.GetDSN())
	}
}

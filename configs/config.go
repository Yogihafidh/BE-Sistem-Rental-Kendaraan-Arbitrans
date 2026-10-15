package configs

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config aggregates all application configuration modules.
type Config struct {
	App   AppConfig
	DB    DatabaseConfig
	JWT   JWTConfig
	MinIO MinIOConfig
	Redis RedisConfig
	CORS  CORSConfig
}

// LoadConfig loads configuration from environment variables and an optional .env file.
// If the .env file does not exist (e.g., in production or container environments),
// it gracefully continues loading from system environment variables.
func LoadConfig(filenames ...string) (*Config, error) {
	// Attempt to load .env file if available
	if len(filenames) > 0 {
		_ = godotenv.Load(filenames...)
	} else {
		// Look for standard .env file; ignore if missing
		if _, err := os.Stat(".env"); err == nil {
			_ = godotenv.Load(".env")
		}
	}

	appEnv := strings.ToLower(getEnv("APP_ENV", "development"))
	appPort := getEnv("APP_PORT", "8080")

	cfg := &Config{
		App: AppConfig{
			Name:            getEnv("APP_NAME", "Arbitrans Vehicle Rental API"),
			Env:             appEnv,
			Port:            appPort,
			ReadTimeout:     getEnvAsDuration("APP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvAsDuration("APP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getEnvAsDuration("APP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getEnvAsDuration("APP_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		DB: DatabaseConfig{
			URL:               getEnv("DATABASE_URL", ""),
			Host:              getEnv("DB_HOST", "localhost"),
			Port:              getEnvAsInt("DB_PORT", 5432),
			User:              getEnv("DB_USER", "postgres"),
			Password:          getEnv("DB_PASSWORD", "postgres"),
			Name:              getEnv("DB_NAME", "rental_db"),
			SSLMode:           getEnv("DB_SSLMODE", "disable"),
			MaxConns:          getEnvAsInt32("DB_MAX_CONNS", 25),
			MinConns:          getEnvAsInt32("DB_MIN_CONNS", 5),
			MaxConnLifetime:   getEnvAsDuration("DB_MAX_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime:   getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", 15*time.Minute),
			HealthCheckPeriod: getEnvAsDuration("DB_HEALTH_CHECK_PERIOD", 1*time.Minute),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", ""),
			ExpirationDuration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
			Issuer:             getEnv("JWT_ISSUER", "arbitrans-rental"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "admin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "password123"),
			Bucket:    getEnv("MINIO_BUCKET", "vehicle-rental"),
			UseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		},
		Redis: RedisConfig{
			Addr:         getEnv("REDIS_ADDR", "localhost:6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			DialTimeout:  getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
			Enabled:      getEnvAsBool("REDIS_ENABLED", true),
		},
		CORS: CORSConfig{
			AllowedOrigins:   getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}, ","),
			AllowedMethods:   getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, ","),
			AllowedHeaders:   getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}, ","),
			ExposeHeaders:    getEnvAsSlice("CORS_EXPOSE_HEADERS", []string{"Content-Length"}, ","),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvAsDuration("CORS_MAX_AGE", 12*time.Hour),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation error: %w", err)
	}

	return cfg, nil
}

// Validate verifies that required configuration values are valid and secure.
func (c *Config) Validate() error {
	var errs []string

	// Validate App Environment
	switch c.App.Env {
	case "development", "staging", "production", "test":
	default:
		errs = append(errs, fmt.Sprintf("invalid APP_ENV '%s': must be development, staging, production, or test", c.App.Env))
	}

	// Validate Port
	portNum, err := strconv.Atoi(c.App.Port)
	if err != nil || portNum < 1 || portNum > 65535 {
		errs = append(errs, fmt.Sprintf("invalid APP_PORT '%s': must be a valid port number (1-65535)", c.App.Port))
	}

	// Validate Database Config
	if strings.TrimSpace(c.DB.URL) == "" {
		if strings.TrimSpace(c.DB.Host) == "" {
			errs = append(errs, "database configuration required: DB_HOST or DATABASE_URL must be specified")
		}
		if strings.TrimSpace(c.DB.Name) == "" {
			errs = append(errs, "database configuration required: DB_NAME must be specified")
		}
		if c.DB.Port < 1 || c.DB.Port > 65535 {
			errs = append(errs, fmt.Sprintf("invalid DB_PORT '%d': must be between 1 and 65535", c.DB.Port))
		}
	}
	if c.DB.MaxConns <= 0 {
		errs = append(errs, "DB_MAX_CONNS must be greater than 0")
	}
	if c.DB.MinConns < 0 || c.DB.MinConns > c.DB.MaxConns {
		errs = append(errs, "DB_MIN_CONNS must be non-negative and less than or equal to DB_MAX_CONNS")
	}

	// Validate JWT Config
	if strings.TrimSpace(c.JWT.Secret) == "" {
		errs = append(errs, "JWT_SECRET is required and cannot be empty")
	} else if c.App.IsProduction() {
		// Production secret security checks
		if len(c.JWT.Secret) < 32 {
			errs = append(errs, "JWT_SECRET must be at least 32 characters in production")
		}
		insecureSecrets := []string{"supersecretkey", "development-secret-change-me", "secret", "admin", "12345678"}
		for _, insecure := range insecureSecrets {
			if strings.EqualFold(c.JWT.Secret, insecure) {
				errs = append(errs, fmt.Sprintf("insecure JWT_SECRET '%s' is prohibited in production", c.JWT.Secret))
				break
			}
		}
	}
	if c.JWT.ExpirationDuration <= 0 {
		errs = append(errs, "JWT_EXPIRATION must be greater than 0")
	}

	// Validate MinIO
	if strings.TrimSpace(c.MinIO.Endpoint) != "" {
		if strings.TrimSpace(c.MinIO.Bucket) == "" {
			errs = append(errs, "MINIO_BUCKET is required when MINIO_ENDPOINT is specified")
		}
		if c.App.IsProduction() {
			if c.MinIO.AccessKey == "admin" || c.MinIO.SecretKey == "password123" {
				errs = append(errs, "default MinIO credentials ('admin' / 'password123') are prohibited in production")
			}
		}
	}

	// Validate Redis
	if c.Redis.Enabled && strings.TrimSpace(c.Redis.Addr) != "" {
		if c.Redis.PoolSize <= 0 {
			errs = append(errs, "REDIS_POOL_SIZE must be greater than 0 when Redis is enabled")
		}
		if c.Redis.DB < 0 {
			errs = append(errs, "REDIS_DB cannot be negative")
		}
	}

	// Validate CORS
	if c.CORS.AllowCredentials && c.App.IsProduction() {
		for _, o := range c.CORS.AllowedOrigins {
			if o == "*" {
				errs = append(errs, "CORS misconfiguration: wildcard origin '*' cannot be combined with AllowCredentials=true in production")
				break
			}
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

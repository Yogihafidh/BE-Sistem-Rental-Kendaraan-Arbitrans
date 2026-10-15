package configs

import (
	"fmt"
	"strings"
	"time"
)

// DatabaseConfig holds PostgreSQL connection pool and credential settings.
type DatabaseConfig struct {
	URL               string
	Host              string
	Port              int
	User              string
	Password          string
	Name              string
	SSLMode           string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

// GetDSN returns the connection string. If URL is provided, it is returned.
// Otherwise, it constructs a PostgreSQL connection string from discrete parameters.
func (d DatabaseConfig) GetDSN() string {
	if strings.TrimSpace(d.URL) != "" {
		return strings.TrimSpace(d.URL)
	}

	sslMode := d.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		sslMode,
	)
}

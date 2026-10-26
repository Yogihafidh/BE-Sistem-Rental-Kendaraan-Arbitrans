package configs

import (
	"fmt"
	"net/url"
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
// Otherwise, it constructs a PostgreSQL connection string from discrete parameters with proper URL encoding.
func (d DatabaseConfig) GetDSN() string {
	if strings.TrimSpace(d.URL) != "" {
		return strings.TrimSpace(d.URL)
	}

	sslMode := d.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:   "/" + d.Name,
	}

	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}

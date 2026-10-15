package configs

import "time"

// AppConfig holds HTTP server and application runtime settings.
type AppConfig struct {
	Name            string
	Env             string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// IsProduction returns true if the application environment is production.
func (a AppConfig) IsProduction() bool {
	return a.Env == "production"
}

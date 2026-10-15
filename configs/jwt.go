package configs

import "time"

// JWTConfig holds JSON Web Token settings.
type JWTConfig struct {
	Secret             string
	ExpirationDuration time.Duration
	Issuer             string
}

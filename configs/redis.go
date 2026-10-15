package configs

import "time"

// RedisConfig holds Redis caching and session store settings.
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Enabled      bool
}

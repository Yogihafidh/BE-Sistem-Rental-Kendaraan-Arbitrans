package redis

import (
	"context"
	"fmt"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/configs"
	goredis "github.com/redis/go-redis/v9"
)

// NewClient initializes a new Redis client with connection pooling and a startup connectivity ping.
// If Redis is disabled in configuration, it returns nil without error.
func NewClient(ctx context.Context, cfg configs.RedisConfig) (*goredis.Client, error) {
	if !cfg.Enabled || cfg.Addr == "" {
		return nil, nil
	}

	opts := &goredis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	client := goredis.NewClient(opts)

	// Verify connectivity with configured dial timeout
	pingCtx, cancel := context.WithTimeout(ctx, cfg.DialTimeout)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to connect to redis at %s: %w", cfg.Addr, err)
	}

	return client, nil
}

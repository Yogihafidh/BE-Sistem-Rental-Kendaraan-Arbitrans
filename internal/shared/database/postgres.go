package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Function untuk membuat PostgreSQL connection pool
func NewPostgresPool(
	ctx context.Context,
) (*pgxpool.Pool, error) {

	// Ambil DATABASE_URL dari environment variable
	dsn := strings.TrimSpace(
		os.Getenv("DATABASE_URL"),
	)

	// Validasi DATABASE_URL
	if dsn == "" {
		return nil, fmt.Errorf(
			"DATABASE_URL is empty",
		)
	}

	// Membuat PostgreSQL connection pool
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	// Return pool connection
	return pool, nil
}

package minio

import (
	"context"
	"fmt"
	"strings"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/configs"
	miniosdk "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps the MinIO client SDK and default bucket.
type Client struct {
	SDK    *miniosdk.Client
	Bucket string
}

// NewClient initializes a MinIO client using the centralized MinIOConfig.
// It returns an error if client initialization fails instead of terminating the process.
func NewClient(cfg configs.MinIOConfig) (*Client, error) {
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("minio endpoint is empty")
	}

	client, err := miniosdk.New(endpoint, &miniosdk.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &Client{
		SDK:    client,
		Bucket: cfg.Bucket,
	}, nil
}

// EnsureBucket verifies that the configured bucket exists, and creates it if it does not.
func (c *Client) EnsureBucket(ctx context.Context) error {
	exists, err := c.SDK.BucketExists(ctx, c.Bucket)
	if err != nil {
		return fmt.Errorf("failed to check minio bucket existence: %w", err)
	}

	if !exists {
		err = c.SDK.MakeBucket(ctx, c.Bucket, miniosdk.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create minio bucket '%s': %w", c.Bucket, err)
		}
	}

	return nil
}

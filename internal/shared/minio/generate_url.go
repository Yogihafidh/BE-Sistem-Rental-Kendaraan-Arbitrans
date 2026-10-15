package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// GetPresignedURL generates a temporary presigned URL for reading an object.
// If expiry is <= 0, a default of 15 minutes is applied.
func (c *Client) GetPresignedURL(
	ctx context.Context,
	objectName string,
	expiry time.Duration,
) (string, error) {
	if expiry <= 0 {
		expiry = 15 * time.Minute
	}

	reqParams := make(url.Values)

	presignedURL, err := c.SDK.PresignedGetObject(
		ctx,
		c.Bucket,
		objectName,
		expiry,
		reqParams,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url for '%s': %w", objectName, err)
	}

	return presignedURL.String(), nil
}

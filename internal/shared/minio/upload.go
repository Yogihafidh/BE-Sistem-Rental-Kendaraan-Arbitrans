package minio

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	miniosdk "github.com/minio/minio-go/v7"
)

// UploadFile uploads a multipart file to MinIO with a unique UUID filename within the given folder.
func (c *Client) UploadFile(
	ctx context.Context,
	file multipart.File,
	header *multipart.FileHeader,
	folder string,
) (string, error) {
	objectName := fmt.Sprintf("%s/%s-%s", folder, uuid.NewString(), header.Filename)

	_, err := c.SDK.PutObject(
		ctx,
		c.Bucket,
		objectName,
		file,
		header.Size,
		miniosdk.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload object to minio: %w", err)
	}

	return objectName, nil
}

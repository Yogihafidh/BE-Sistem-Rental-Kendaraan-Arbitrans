package storage

import (
	"context"
	"net/url"
	"time"
)

func (m *MinioClient) GetPresignedURL(
	objectName string,
) (string, error) {
	// Membuat query parameter kosong untuk presigned URL
	reqParams := make(url.Values)

	// Membuat temporary signed URL untuk akses file object di MinIO dengan masa berlaku 15 menit
	presignedURL, err := m.Client.PresignedGetObject(
		context.Background(),
		m.Bucket,
		objectName,
		time.Minute*15,
		reqParams,
	)

	// Jika gagal generate URL
	if err != nil {
		return "", err
	}

	// Convert URL object menjadi string
	return presignedURL.String(), nil
}

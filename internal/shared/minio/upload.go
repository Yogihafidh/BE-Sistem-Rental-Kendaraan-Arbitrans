package storage

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"       // Library untuk generate UUID random
	"github.com/minio/minio-go/v7" // Library client MinIO
)

func (m *MinioClient) UploadFile(
	file multipart.File,
	header *multipart.FileHeader,
	folder string,
) (string, error) {

	// Membuat nama object/file yang unik dengan format: folder/uuid-random-nama-file-asli
	objectName := folder + "/" + uuid.NewString() + "-" + header.Filename

	// Upload file ke MinIO bucket dengan nama object yang sudah dibuat
	_, err := m.Client.PutObject(
		context.Background(),
		m.Bucket,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)

	// Jika upload gagal
	if err != nil {
		return "", err
	}

	// Return nama/path object file
	return objectName, nil
}

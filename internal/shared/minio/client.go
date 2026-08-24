package storage

import (
	"log"

	// Library MinIO client
	"github.com/minio/minio-go/v7"

	// Package credentials untuk autentikasi MinIO
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client // Object client MinIO utama
	Bucket string        // Nama bucket default yang akan dipakai aplikasi
}

func NewMinioClient() *MinioClient {
	// Endpoint MinIO Karena MinIO berjalan di Docker lokal pada port 9000
	endpoint := "localhost:9000"

	// Username MinIO
	accessKey := "admin"

	// Password MinIO
	secretKey := "password123"

	// Membuat koneksi client MinIO
	client, err := minio.New(endpoint, &minio.Options{
		// Konfigurasi credentials login
		Creds: credentials.NewStaticV4(
			accessKey,
			secretKey,
			"",
		),

		// false artinya tidak menggunakan HTTPS karena local development
		Secure: false,
	})

	// Jika gagal koneksi ke MinIO
	if err != nil {
		log.Fatal(err)
	}

	// Return object MinioClient dengan koneksi client yang sudah dibuat dan bucket default
	return &MinioClient{
		Client: client,
		Bucket: "vehicle-rental",
	}
}

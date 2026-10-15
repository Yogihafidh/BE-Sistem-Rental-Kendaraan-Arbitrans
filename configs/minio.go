package configs

// MinIOConfig holds S3/MinIO object storage settings.
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

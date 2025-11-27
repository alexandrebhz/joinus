package storage

import (
	"fmt"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/infrastructure/config"
)

func NewStorageService(cfg *config.Config) (port.StorageService, error) {
	if cfg.Storage.Type == "s3" {
		// Validate S3 config
		if cfg.Storage.S3.Bucket == "" || cfg.Storage.S3.AccessKey == "" || cfg.Storage.S3.SecretKey == "" {
			return nil, fmt.Errorf("S3 storage requires bucket, access key, and secret key")
		}
		return NewS3Storage(cfg.Storage.S3)
	}
	
	// For MinIO, check if endpoint is accessible
	if cfg.Storage.MinIO.Endpoint == "" {
		return nil, fmt.Errorf("MinIO endpoint is not configured")
	}
	
	return NewMinIOStorage(cfg.Storage.MinIO)
}




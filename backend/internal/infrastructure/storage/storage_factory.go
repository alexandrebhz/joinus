package storage

import (
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/infrastructure/config"
)

func NewStorageService(cfg *config.Config) (port.StorageService, error) {
	if cfg.Storage.Type == "s3" {
		return NewS3Storage(cfg.Storage.S3)
	}
	return NewMinIOStorage(cfg.Storage.MinIO)
}




package storage

import (
	"context"
	"bytes"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/infrastructure/config"
)

type MinIOStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinIOStorage(cfg config.MinIOConfig) (port.StorageService, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	storage := &MinIOStorage{
		client:     client,
		bucketName: cfg.Bucket,
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return storage, nil
}

func (s *MinIOStorage) Upload(ctx context.Context, fileData []byte, key string, contentType string) (string, error) {
	reader := bytes.NewReader(fileData)
	
	_, err := s.client.PutObject(ctx, s.bucketName, key, reader, int64(len(fileData)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return URL (in production, this would be a proper URL)
	url := fmt.Sprintf("http://%s/%s/%s", s.client.EndpointURL().Host, s.bucketName, key)
	return url, nil
}

func (s *MinIOStorage) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{})
}

func (s *MinIOStorage) GetURL(ctx context.Context, key string) (string, error) {
	url := fmt.Sprintf("http://%s/%s/%s", s.client.EndpointURL().Host, s.bucketName, key)
	return url, nil
}




package port

import "context"

type StorageService interface {
	Upload(ctx context.Context, file []byte, key string, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
}




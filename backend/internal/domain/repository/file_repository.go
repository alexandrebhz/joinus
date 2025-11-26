package repository

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type FileRepository interface {
	Create(ctx context.Context, file *entity.File) error
	FindByID(ctx context.Context, id string) (*entity.File, error)
	Delete(ctx context.Context, id string) error
}




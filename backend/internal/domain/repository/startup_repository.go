package repository

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type Pagination struct {
	Page     int
	PageSize int
	OrderBy  string
	OrderDir string // ASC, DESC
}

type StartupRepository interface {
	Create(ctx context.Context, startup *entity.Startup) error
	Update(ctx context.Context, startup *entity.Startup) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.Startup, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Startup, error)
	FindByAPIToken(ctx context.Context, token string) (*entity.Startup, error)
	List(ctx context.Context, filter StartupFilter) ([]*entity.Startup, int64, error)
}

type StartupFilter struct {
	Industry string
	Status   entity.StartupStatus
	Search   string
	Pagination
}




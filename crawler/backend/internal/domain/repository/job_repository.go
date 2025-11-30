package repository

import (
	"context"
	"github.com/startup-job-board/crawler/internal/domain/entity"
)

// JobRepository defines the interface for crawled job persistence
type JobRepository interface {
	Create(ctx context.Context, job *entity.CrawledJob) error
	Update(ctx context.Context, job *entity.CrawledJob) error
	FindByID(ctx context.Context, id string) (*entity.CrawledJob, error)
	FindBySiteID(ctx context.Context, siteID string) ([]*entity.CrawledJob, error)
	FindUnsynced(ctx context.Context, limit int) ([]*entity.CrawledJob, error)
	ExistsByURL(ctx context.Context, url string) (bool, error)
	ExistsByHash(ctx context.Context, hash string) (bool, error)
	ExistsByExternalID(ctx context.Context, externalID string) (bool, error)
	MarkAsSynced(ctx context.Context, id string) error
	CountUnsynced(ctx context.Context) (int64, error)
}


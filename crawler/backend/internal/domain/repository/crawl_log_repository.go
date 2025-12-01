package repository

import (
	"context"
	"github.com/startup-job-board/crawler/internal/domain/entity"
)

// CrawlLogRepository defines the interface for crawl log persistence
type CrawlLogRepository interface {
	Create(ctx context.Context, log *entity.CrawlLog) error
	Update(ctx context.Context, log *entity.CrawlLog) error
	FindByID(ctx context.Context, id string) (*entity.CrawlLog, error)
	FindBySiteID(ctx context.Context, siteID string, limit int) ([]*entity.CrawlLog, error)
	FindLatestBySiteID(ctx context.Context, siteID string) (*entity.CrawlLog, error)
}


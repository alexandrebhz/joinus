package repository

import (
	"context"
	"github.com/startup-job-board/crawler/internal/domain/entity"
)

// SiteRepository defines the interface for crawl site persistence
type SiteRepository interface {
	Create(ctx context.Context, site *entity.CrawlSite) error
	Update(ctx context.Context, site *entity.CrawlSite) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.CrawlSite, error)
	FindAll(ctx context.Context) ([]*entity.CrawlSite, error)
	FindActive(ctx context.Context) ([]*entity.CrawlSite, error)
	UpdateLastCrawledAt(ctx context.Context, id string) error
}


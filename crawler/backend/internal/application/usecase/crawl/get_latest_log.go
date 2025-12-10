package crawl

import (
	"context"

	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// GetLatestLogUseCase handles retrieving the latest crawl log
type GetLatestLogUseCase struct {
	crawlLogRepo repository.CrawlLogRepository
}

// NewGetLatestLogUseCase creates a new GetLatestLogUseCase
func NewGetLatestLogUseCase(crawlLogRepo repository.CrawlLogRepository) *GetLatestLogUseCase {
	return &GetLatestLogUseCase{
		crawlLogRepo: crawlLogRepo,
	}
}

// Execute retrieves the latest crawl log for a site
func (uc *GetLatestLogUseCase) Execute(ctx context.Context, siteID string) (*entity.CrawlLog, error) {
	return uc.crawlLogRepo.FindLatestBySiteID(ctx, siteID)
}







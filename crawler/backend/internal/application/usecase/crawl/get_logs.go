package crawl

import (
	"context"

	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// GetLogsUseCase handles retrieving crawl logs
type GetLogsUseCase struct {
	crawlLogRepo repository.CrawlLogRepository
}

// NewGetLogsUseCase creates a new GetLogsUseCase
func NewGetLogsUseCase(crawlLogRepo repository.CrawlLogRepository) *GetLogsUseCase {
	return &GetLogsUseCase{
		crawlLogRepo: crawlLogRepo,
	}
}

// Execute retrieves crawl logs for a site
func (uc *GetLogsUseCase) Execute(ctx context.Context, siteID string, limit int) ([]*entity.CrawlLog, error) {
	return uc.crawlLogRepo.FindBySiteID(ctx, siteID, limit)
}




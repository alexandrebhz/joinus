package crawl

import (
	"context"

	"github.com/startup-job-board/crawler/internal/domain/repository"
	"github.com/startup-job-board/crawler/internal/infrastructure/crawler"
)

// ExecuteCrawlUseCase handles manual crawl execution
type ExecuteCrawlUseCase struct {
	crawler  *crawler.Engine
	siteRepo repository.SiteRepository
}

// NewExecuteCrawlUseCase creates a new ExecuteCrawlUseCase
func NewExecuteCrawlUseCase(crawler *crawler.Engine, siteRepo repository.SiteRepository) *ExecuteCrawlUseCase {
	return &ExecuteCrawlUseCase{
		crawler:  crawler,
		siteRepo: siteRepo,
	}
}

// Execute executes a crawl for a specific site
func (uc *ExecuteCrawlUseCase) Execute(ctx context.Context, siteID string) (*crawler.CrawlResult, error) {
	site, err := uc.siteRepo.FindByID(ctx, siteID)
	if err != nil {
		return nil, err
	}

	result, err := uc.crawler.Crawl(ctx, site)
	if err != nil {
		return nil, err
	}

	// Update last crawled timestamp
	uc.siteRepo.UpdateLastCrawledAt(ctx, siteID)

	return result, nil
}


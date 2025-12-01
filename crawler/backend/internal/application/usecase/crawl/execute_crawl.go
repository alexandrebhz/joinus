package crawl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
	"github.com/startup-job-board/crawler/internal/infrastructure/crawler"
)

// ExecuteCrawlUseCase handles manual crawl execution
type ExecuteCrawlUseCase struct {
	crawler     *crawler.Engine
	siteRepo    repository.SiteRepository
	crawlLogRepo repository.CrawlLogRepository
}

// NewExecuteCrawlUseCase creates a new ExecuteCrawlUseCase
func NewExecuteCrawlUseCase(crawler *crawler.Engine, siteRepo repository.SiteRepository, crawlLogRepo repository.CrawlLogRepository) *ExecuteCrawlUseCase {
	return &ExecuteCrawlUseCase{
		crawler:     crawler,
		siteRepo:    siteRepo,
		crawlLogRepo: crawlLogRepo,
	}
}

// Execute executes a crawl for a specific site
func (uc *ExecuteCrawlUseCase) Execute(ctx context.Context, siteID string) (*crawler.CrawlResult, error) {
	site, err := uc.siteRepo.FindByID(ctx, siteID)
	if err != nil {
		return nil, err
	}

	// Create crawl log
	crawlLog := &entity.CrawlLog{
		ID:        uuid.New().String(),
		SiteID:    siteID,
		Status:    "running",
		StartedAt: time.Now(),
		Errors:    []string{},
		Logs:      []entity.LogEntry{},
	}
	crawlLog.AddLog("info", "Starting crawl for site: "+site.Name)
	crawlLog.AddLog("info", "Base URL: "+site.BaseURL)

	if err := uc.crawlLogRepo.Create(ctx, crawlLog); err != nil {
		// Log error but continue with crawl
		crawlLog.AddLog("warning", "Failed to create crawl log: "+err.Error())
	}

	// Execute crawl
	crawlLog.AddLog("info", "Fetching pages...")
	result, err := uc.crawler.Crawl(ctx, site)
	
	if err != nil {
		crawlLog.Fail(err)
		uc.crawlLogRepo.Update(ctx, crawlLog)
		return nil, err
	}
	
	// Update log with results
	crawlLog.JobsFound = result.JobsFound
	crawlLog.JobsSaved = result.JobsSaved
	crawlLog.JobsSkipped = result.JobsSkipped
	crawlLog.PagesCrawled = result.PagesCrawled
	crawlLog.Errors = result.Errors
	
	// Add detailed logs for errors
	if len(result.Errors) > 0 {
		for _, errMsg := range result.Errors {
			crawlLog.AddLog("error", errMsg)
		}
	}

	// Complete log
	crawlLog.Complete()
	crawlLog.AddLog("info", fmt.Sprintf("Crawl completed: %d jobs found, %d saved, %d skipped", result.JobsFound, result.JobsSaved, result.JobsSkipped))
	uc.crawlLogRepo.Update(ctx, crawlLog)

	// Update last crawled timestamp
	uc.siteRepo.UpdateLastCrawledAt(ctx, siteID)

	return result, nil
}


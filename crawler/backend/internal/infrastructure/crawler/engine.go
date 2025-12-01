package crawler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// Engine handles the crawling process
type Engine struct {
	jobRepo repository.JobRepository
	client  *http.Client
}

// NewEngine creates a new crawler engine
func NewEngine(jobRepo repository.JobRepository) *Engine {
	return &Engine{
		jobRepo: jobRepo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CrawlResult represents the result of a crawl operation
type CrawlResult struct {
	JobsFound    int      `json:"jobs_found"`
	JobsSaved    int      `json:"jobs_saved"`
	JobsSkipped  int      `json:"jobs_skipped"`
	PagesCrawled int      `json:"pages_crawled"`
	Errors       []string `json:"errors"`
}

// Crawl crawls a site for job listings
func (e *Engine) Crawl(ctx context.Context, site *entity.CrawlSite) (*CrawlResult, error) {
	result := &CrawlResult{
		Errors:       []string{},
		PagesCrawled: 0,
	}

	extractor := NewExtractor(site.ExtractionRules)
	paginator := NewPaginator(site.PaginationConfig)

	// Get page URLs to crawl
	pageURLs, err := paginator.GetPageURLs(site.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate page URLs: %w", err)
	}

	// Crawl each page
	for _, pageURL := range pageURLs {
		result.PagesCrawled++
		// Respect request delay
		if site.RequestDelay > 0 {
			time.Sleep(time.Duration(site.RequestDelay) * time.Second)
		}

		// Fetch page
		doc, err := e.fetchPage(ctx, pageURL, site.UserAgent)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to fetch %s: %v", pageURL, err))
			continue
		}

		// Extract jobs
		jobs, err := extractor.ExtractJobs(doc, pageURL)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to extract jobs from %s: %v", pageURL, err))
			continue
		}

		result.JobsFound += len(jobs)

		// Save jobs
		for _, job := range jobs {
			if err := e.saveJob(ctx, job, site); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to save job: %v", err))
				result.JobsSkipped++
				continue
			}
			result.JobsSaved++
		}

		// Check for next page (for link_follow pagination)
		if paginator.HasNextPage(doc) {
			nextURL, err := paginator.GetNextPageURL(doc, pageURL)
			if err == nil && nextURL != "" {
				pageURLs = append(pageURLs, nextURL)
			}
		}
	}

	return result, nil
}

// fetchPage fetches a page and returns a goquery document
func (e *Engine) fetchPage(ctx context.Context, url string, userAgent string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; JobCrawler/1.0)")
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// saveJob saves a job if it's not a duplicate
func (e *Engine) saveJob(ctx context.Context, job *entity.CrawledJob, site *entity.CrawlSite) error {
	// Set site ID
	job.SiteID = site.ID

	// Generate deduplication hash
	job.DeduplicationHash = job.GenerateDeduplicationHash(site.DeduplicationKey)

	// Check for duplicates
	var exists bool
	var err error

	switch site.DeduplicationKey {
	case "url":
		exists, err = e.jobRepo.ExistsByURL(ctx, job.DetailURL)
	case "composite":
		exists, err = e.jobRepo.ExistsByHash(ctx, job.DeduplicationHash)
	case "external_id":
		if job.ExternalID != "" {
			exists, err = e.jobRepo.ExistsByExternalID(ctx, job.ExternalID)
		}
	default:
		exists, err = e.jobRepo.ExistsByURL(ctx, job.DetailURL)
	}

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("job already exists")
	}

	// Generate ID and timestamps
	job.ID = uuid.New().String()
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	job.Synced = false

	// Save job
	return e.jobRepo.Create(ctx, job)
}

// SetHTTPClient sets a custom HTTP client (useful for testing)
func (e *Engine) SetHTTPClient(client *http.Client) {
	e.client = client
}


# Job Crawler System - Technical Specification

## Overview

The Job Crawler System is a distributed web scraping service designed to aggregate job listings from multiple external sources and synchronize them with the JoinUs backend API. The system provides a flexible, configurable interface for managing crawl targets, extraction rules, and scheduling policies.

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Crawler Service                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│  │ Site Manager │───▶│   Scheduler  │───▶│   Crawler    │ │
│  │  (CRUD API)  │    │  (Cron Jobs) │    │   Engine     │ │
│  └──────────────┘    └──────────────┘    └──────────────┘ │
│         │                    │                    │         │
│         └────────────────────┼────────────────────┘         │
│                              │                              │
│                    ┌─────────▼─────────┐                   │
│                    │   Local Database  │                   │
│                    │   (PostgreSQL)    │                   │
│                    └─────────┬─────────┘                   │
│                              │                              │
│                    ┌─────────▼─────────┐                   │
│                    │   Sync Service    │                   │
│                    │  (Backend API)    │                   │
│                    └───────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

### Technology Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL (local storage for crawled jobs)
- **HTTP Client**: Standard `net/http` or `colly` framework
- **HTML Parsing**: `goquery` or `colly`
- **Scheduling**: `robfig/cron` v3
- **API Client**: HTTP client for backend synchronization

## Core Features

### 1. Site Configuration Management (CRUD Interface)

Each crawl target (site) requires configuration for:
- **Base URL**: The entry point URL containing job listings
- **Pagination Strategy**: How to navigate through multiple pages
- **Field Extraction Rules**: CSS/XPath selectors for each job field
- **Crawl Schedule**: Cron expression for automated crawling
- **Deduplication Strategy**: How to identify duplicate jobs

#### Site Configuration Schema

```go
type CrawlSite struct {
    ID                  string    `json:"id"`
    Name                string    `json:"name"`                // Human-readable name
    BaseURL             string    `json:"base_url"`            // Entry point URL
    Active              bool      `json:"active"`              // Enable/disable crawling
    Schedule            string    `json:"schedule"`            // Cron expression (e.g., "0 0 * * *")
    LastCrawledAt       *time.Time `json:"last_crawled_at"`
    NextCrawlAt         *time.Time `json:"next_crawl_at"`
    CrawlInterval       string    `json:"crawl_interval"`      // "daily", "weekly", "custom"
    PaginationConfig     PaginationConfig `json:"pagination_config"`
    ExtractionRules     ExtractionRules   `json:"extraction_rules"`
    DeduplicationKey    string    `json:"deduplication_key"`   // Field to use for deduplication
    CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
}
```

### 2. Pagination Configuration

Support multiple pagination strategies:

#### Strategy Types

**A. Query Parameter Pagination**
```json
{
  "type": "query_param",
  "param_name": "page",
  "start_page": 1,
  "increment": 1,
  "max_pages": 100,
  "next_page_selector": ".pagination .next"
}
```

**B. URL Pattern Pagination**
```json
{
  "type": "url_pattern",
  "pattern": "/jobs/page/{page}",
  "start_page": 1,
  "increment": 1,
  "max_pages": 50
}
```

**C. Link Following Pagination**
```json
{
  "type": "link_follow",
  "next_link_selector": "a.pagination-next",
  "max_depth": 10
}
```

**D. Infinite Scroll / API Pagination**
```json
{
  "type": "api_pagination",
  "api_endpoint": "https://example.com/api/jobs",
  "page_param": "offset",
  "page_size": 20,
  "max_pages": 100
}
```

#### Pagination Config Schema

```go
type PaginationConfig struct {
    Type            string `json:"type"`              // "query_param", "url_pattern", "link_follow", "api_pagination"
    ParamName       string `json:"param_name,omitempty"`
    StartPage       int    `json:"start_page,omitempty"`
    Increment       int    `json:"increment,omitempty"`
    MaxPages        int    `json:"max_pages,omitempty"`
    NextPageSelector string `json:"next_page_selector,omitempty"`
    URLPattern      string `json:"url_pattern,omitempty"`
    APIConfig       *APIPaginationConfig `json:"api_config,omitempty"`
}
```

### 3. Field Extraction Rules

Each job field requires a CSS selector or XPath expression to extract data from HTML.

#### Extraction Rules Schema

```go
type ExtractionRules struct {
    JobListSelector    string            `json:"job_list_selector"`    // Selector for job container (e.g., ".job-item")
    JobDetailURL       JobURLRule        `json:"job_detail_url"`       // How to get job detail page URL
    Fields             map[string]FieldRule `json:"fields"`            // Field-specific extraction rules
}

type FieldRule struct {
    Selector          string   `json:"selector"`           // CSS selector or XPath
    Type              string   `json:"type"`               // "text", "html", "attribute", "regex"
    Attribute         string   `json:"attribute,omitempty"` // For attribute extraction (e.g., "href", "data-id")
    RegexPattern      string   `json:"regex_pattern,omitempty"` // For regex extraction
    Required          bool     `json:"required"`          // Whether field is mandatory
    DefaultValue      string   `json:"default_value,omitempty"`
    Transformations   []string `json:"transformations,omitempty"` // e.g., ["trim", "lowercase", "parse_date"]
}

type JobURLRule struct {
    Type              string `json:"type"`                // "relative", "absolute", "attribute"
    Selector          string `json:"selector"`            // Selector for link element
    Attribute         string `json:"attribute,omitempty"` // Usually "href"
    BaseURL           string `json:"base_url,omitempty"`  // For relative URLs
}
```

#### Example Extraction Rules

```json
{
  "job_list_selector": ".job-listing-item",
  "job_detail_url": {
    "type": "relative",
    "selector": "a.job-link",
    "attribute": "href",
    "base_url": "https://example.com"
  },
  "fields": {
    "title": {
      "selector": "h2.job-title",
      "type": "text",
      "required": true
    },
    "description": {
      "selector": ".job-description",
      "type": "html",
      "required": true,
      "transformations": ["strip_html", "trim"]
    },
    "company": {
      "selector": ".company-name",
      "type": "text",
      "required": false
    },
    "location": {
      "selector": ".job-location",
      "type": "text",
      "required": false
    },
    "salary_min": {
      "selector": ".salary-range",
      "type": "regex",
      "regex_pattern": "\\$([0-9,]+)",
      "required": false,
      "transformations": ["remove_commas", "parse_int"]
    },
    "application_url": {
      "selector": "a.apply-button",
      "type": "attribute",
      "attribute": "href",
      "required": false
    }
  }
}
```

### 4. Deduplication Strategy

Prevent duplicate job entries using configurable deduplication keys.

#### Deduplication Methods

**A. URL-based Deduplication** (Recommended)
- Use job detail URL as unique identifier
- Most reliable for detecting duplicates

**B. Composite Key Deduplication**
- Combine multiple fields (e.g., title + company + location)
- Use hash of combined fields as key

**C. External ID Deduplication**
- Use site-specific job ID if available
- Extract from HTML attribute or URL pattern

#### Deduplication Check Flow

```go
func (c *Crawler) isDuplicate(job *CrawledJob) (bool, error) {
    switch c.site.DeduplicationKey {
    case "url":
        return c.repo.ExistsByURL(job.DetailURL)
    case "composite":
        hash := generateHash(job.Title, job.Company, job.Location)
        return c.repo.ExistsByHash(hash)
    case "external_id":
        return c.repo.ExistsByExternalID(job.ExternalID)
    default:
        return false, nil
    }
}
```

### 5. Scheduling System

Use cron-style scheduling with configurable intervals per site.

#### Schedule Configuration

```go
type ScheduleConfig struct {
    CronExpression string `json:"cron_expression"` // Standard cron format
    Timezone       string `json:"timezone"`         // e.g., "America/New_York"
    Enabled        bool   `json:"enabled"`
}
```

#### Common Schedule Patterns

- **Daily at midnight**: `0 0 * * *`
- **Every 6 hours**: `0 */6 * * *`
- **Weekly (Monday 9 AM)**: `0 9 * * 1`
- **Every 30 minutes**: `*/30 * * * *`

#### Scheduler Implementation

```go
type Scheduler struct {
    cron      *cron.Cron
    crawler   *CrawlerEngine
    siteRepo  SiteRepository
}

func (s *Scheduler) Start() error {
    sites, err := s.siteRepo.FindActive()
    if err != nil {
        return err
    }
    
    for _, site := range sites {
        s.scheduleSite(site)
    }
    
    s.cron.Start()
    return nil
}

func (s *Scheduler) scheduleSite(site *CrawlSite) {
    s.cron.AddFunc(site.Schedule, func() {
        s.crawler.Crawl(site)
    })
}
```

### 6. Local Database Schema

Store crawled jobs locally before synchronization.

#### Crawled Job Entity

```go
type CrawledJob struct {
    ID                string    `json:"id"`                 // UUID
    SiteID            string    `json:"site_id"`            // Reference to CrawlSite
    ExternalID        string    `json:"external_id"`        // Site-specific job ID
    DetailURL         string    `json:"detail_url"`         // Full URL to job detail page
    Title             string    `json:"title"`
    Description       string    `json:"description"`
    Requirements      string    `json:"requirements"`
    Company           string    `json:"company"`
    Location          string    `json:"location"`
    City              string    `json:"city"`
    Country           string    `json:"country"`
    JobType           string    `json:"job_type"`           // "full_time", "part_time", etc.
    LocationType      string    `json:"location_type"`     // "remote", "hybrid", "onsite"
    SalaryMin         *int      `json:"salary_min"`
    SalaryMax         *int      `json:"salary_max"`
    Currency          string    `json:"currency"`
    ApplicationURL    *string   `json:"application_url"`
    ApplicationEmail  *string   `json:"application_email"`
    ExpiresAt         *time.Time `json:"expires_at"`
    RawHTML           string    `json:"raw_html"`           // Store raw HTML for debugging
    DeduplicationHash string   `json:"deduplication_hash"` // For duplicate detection
    Synced            bool      `json:"synced"`             // Whether synced to backend
    SyncedAt          *time.Time `json:"synced_at"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}
```

### 7. Backend API Synchronization

Periodically sync crawled jobs to the backend API.

#### Sync Strategy

**A. Batch Sync**
- Collect unsynced jobs
- Send in batches (e.g., 50 jobs per request)
- Handle partial failures gracefully

**B. Incremental Sync**
- Sync only new jobs since last sync
- Use timestamp-based filtering

**C. Full Sync (On-demand)**
- Manual trigger for full re-sync
- Useful for data recovery or migration

#### Sync Service Interface

```go
type SyncService interface {
    SyncJobs(ctx context.Context, jobs []*CrawledJob) (*SyncResult, error)
    SyncJob(ctx context.Context, job *CrawledJob) error
    GetBackendStartupID(ctx context.Context, siteID string) (string, error)
}

type SyncResult struct {
    SuccessCount int      `json:"success_count"`
    FailureCount int      `json:"failure_count"`
    Errors       []error  `json:"errors"`
}
```

#### Backend API Integration

The crawler must map crawled jobs to backend job DTOs:

```go
func (c *Crawler) mapToBackendJob(crawled *CrawledJob, startupID string) (*dto.CreateJobInput, error) {
    return &dto.CreateJobInput{
        StartupID:      startupID,
        Title:          crawled.Title,
        Description:    crawled.Description,
        Requirements:   crawled.Requirements,
        JobType:        normalizeJobType(crawled.JobType),
        LocationType:   normalizeLocationType(crawled.LocationType),
        City:           crawled.City,
        Country:        crawled.Country,
        SalaryMin:      crawled.SalaryMin,
        SalaryMax:      crawled.SalaryMax,
        Currency:       normalizeCurrency(crawled.Currency),
        ApplicationURL: crawled.ApplicationURL,
        ApplicationEmail: crawled.ApplicationEmail,
        ExpiresAt:      formatExpiresAt(crawled.ExpiresAt),
    }, nil
}
```

## Implementation Considerations

### Error Handling

- **Network Errors**: Retry with exponential backoff
- **Parsing Errors**: Log and continue with other jobs
- **Rate Limiting**: Respect `robots.txt` and implement delays
- **HTML Structure Changes**: Alert when extraction fails consistently

### Rate Limiting & Ethics

- Respect `robots.txt` files
- Implement delays between requests (configurable per site)
- Use proper User-Agent headers
- Monitor for IP bans or CAPTCHAs

### Monitoring & Observability

- **Metrics**: Jobs crawled, sync success rate, errors per site
- **Logging**: Structured logging for debugging
- **Health Checks**: Scheduler status, database connectivity
- **Alerts**: Failed crawls, sync failures, site structure changes

### Data Quality

- **Validation**: Validate extracted data against backend DTOs
- **Normalization**: Standardize job types, location types, currencies
- **Sanitization**: Clean HTML, remove special characters
- **Enrichment**: Geocode locations, parse salary ranges

## Questions & Open Decisions

1. **Authentication**: How should the crawler authenticate with the backend API?
   - API Token per site?
   - Single service account?
   - OAuth2 client credentials?

2. **Startup Mapping**: How are external sites mapped to backend startups?
   - Manual configuration (site → startup_id)?
   - Automatic matching by company name?
   - Separate mapping table?

3. **Job Ownership**: Should crawled jobs be:
   - Associated with a specific startup entity?
   - Created as "unclaimed" jobs?
   - Require manual approval before publishing?

4. **Field Mapping**: How to handle sites with different field structures?
   - Manual field mapping configuration?
   - AI-powered field detection?
   - Template-based extraction?

5. **Failure Handling**: What happens when:
   - Backend API is unavailable during sync?
   - Job extraction fails for a site?
   - Site structure changes and selectors break?

6. **Storage**: Should raw HTML be stored?
   - For debugging purposes?
   - For re-extraction if rules change?
   - Storage limits/retention policies?

7. **Concurrency**: How many sites should be crawled concurrently?
   - Per-site rate limiting?
   - Global rate limiting?
   - Queue-based processing?

8. **Notification**: How to notify about:
   - New jobs discovered?
   - Crawl failures?
   - Sync errors?
   - Site structure changes?

## Proposed Project Structure

```
crawler/
├── cmd/
│   └── crawler/
│       └── main.go                 # Entry point
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── crawl_site.go
│   │   │   └── crawled_job.go
│   │   └── repository/
│   │       ├── site_repository.go
│   │       └── job_repository.go
│   ├── application/
│   │   ├── usecase/
│   │   │   ├── site/
│   │   │   │   ├── create_site.go
│   │   │   │   ├── update_site.go
│   │   │   │   └── delete_site.go
│   │   │   └── crawl/
│   │   │       ├── execute_crawl.go
│   │   │       └── sync_jobs.go
│   │   └── dto/
│   │       ├── site_dto.go
│   │       └── job_dto.go
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   └── postgres/
│   │   │       ├── site_repository_impl.go
│   │   │       └── job_repository_impl.go
│   │   ├── crawler/
│   │   │   ├── engine.go
│   │   │   ├── extractor.go
│   │   │   └── paginator.go
│   │   ├── scheduler/
│   │   │   └── cron_scheduler.go
│   │   └── sync/
│   │       └── api_sync_service.go
│   └── presentation/
│       └── http/
│           ├── handler/
│           │   └── site_handler.go
│           └── router/
│               └── router.go
├── pkg/
│   ├── html/
│   │   └── selector.go
│   ├── normalization/
│   │   ├── job_type.go
│   │   └── location.go
│   └── logger/
│       └── logger.go
├── migrations/
│   └── 001_initial_schema.sql
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## Next Steps

1. **Clarify open questions** with stakeholders
2. **Design database schema** for sites and crawled jobs
3. **Implement site CRUD API** for configuration management
4. **Build crawler engine** with extraction and pagination
5. **Implement scheduler** with cron integration
6. **Create sync service** for backend API integration
7. **Add monitoring and logging** infrastructure
8. **Write tests** for critical components

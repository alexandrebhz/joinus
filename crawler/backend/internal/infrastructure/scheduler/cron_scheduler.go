package scheduler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
	"github.com/startup-job-board/crawler/internal/infrastructure/crawler"
)

// Scheduler manages cron jobs for crawling sites
type Scheduler struct {
	cron     *cron.Cron
	crawler  *crawler.Engine
	siteRepo repository.SiteRepository
	entries  map[string]cron.EntryID
}

// NewScheduler creates a new scheduler
func NewScheduler(crawler *crawler.Engine, siteRepo repository.SiteRepository) *Scheduler {
	c := cron.New(cron.WithSeconds()) // Use seconds for more flexibility

	return &Scheduler{
		cron:     c,
		crawler:  crawler,
		siteRepo: siteRepo,
		entries:  make(map[string]cron.EntryID),
	}
}

// Start starts the scheduler and loads all active sites
func (s *Scheduler) Start(ctx context.Context) error {
	sites, err := s.siteRepo.FindActive(ctx)
	if err != nil {
		return err
	}

	for _, site := range sites {
		s.ScheduleSite(ctx, site)
	}

	s.cron.Start()
	log.Println("Scheduler started")
	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// ScheduleSite schedules a site for crawling
func (s *Scheduler) ScheduleSite(ctx context.Context, site *entity.CrawlSite) {
	// Remove existing entry if any
	if entryID, exists := s.entries[site.ID]; exists {
		s.cron.Remove(entryID)
	}

	// Add new cron job
	entryID, err := s.cron.AddFunc(site.Schedule, func() {
		log.Printf("Starting crawl for site: %s (%s)", site.Name, site.ID)
		
		result, err := s.crawler.Crawl(ctx, site)
		if err != nil {
			log.Printf("Error crawling site %s: %v", site.Name, err)
			return
		}

		log.Printf("Crawl completed for site %s: %d jobs found, %d saved, %d skipped",
			site.Name, result.JobsFound, result.JobsSaved, result.JobsSkipped)

		// Update last crawled timestamp
		s.siteRepo.UpdateLastCrawledAt(ctx, site.ID)
	})

	if err != nil {
		log.Printf("Error scheduling site %s: %v", site.Name, err)
		return
	}

	s.entries[site.ID] = entryID
	log.Printf("Scheduled site %s with schedule: %s", site.Name, site.Schedule)
}

// UnscheduleSite removes a site from the scheduler
func (s *Scheduler) UnscheduleSite(siteID string) {
	if entryID, exists := s.entries[siteID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, siteID)
		log.Printf("Unscheduled site: %s", siteID)
	}
}

// ReloadSite reloads a site's schedule
func (s *Scheduler) ReloadSite(ctx context.Context, siteID string) error {
	site, err := s.siteRepo.FindByID(ctx, siteID)
	if err != nil {
		return err
	}

	if site.Active {
		s.ScheduleSite(ctx, site)
	} else {
		s.UnscheduleSite(siteID)
	}

	return nil
}


package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/startup-job-board/crawler/internal/application/usecase/crawl"
	"github.com/startup-job-board/crawler/internal/application/usecase/site"
	"github.com/startup-job-board/crawler/internal/config"
	"github.com/startup-job-board/crawler/internal/infrastructure/crawler"
	"github.com/startup-job-board/crawler/internal/infrastructure/persistence/postgres"
	"github.com/startup-job-board/crawler/internal/infrastructure/scheduler"
	"github.com/startup-job-board/crawler/internal/infrastructure/sync"
	"github.com/startup-job-board/crawler/internal/presentation/http/handler"
	"github.com/startup-job-board/crawler/internal/presentation/http/router"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	siteRepo := postgres.NewSiteRepository(db)
	jobRepo := postgres.NewJobRepository(db)
	crawlLogRepo := postgres.NewCrawlLogRepository(db)

	// Initialize crawler engine
	crawlerEngine := crawler.NewEngine(jobRepo)

	// Initialize scheduler
	scheduler := scheduler.NewScheduler(crawlerEngine, siteRepo)

	// Initialize sync service
	syncService := sync.NewSyncService(cfg.BackendURL, cfg.BackendToken, jobRepo, siteRepo)

	// Initialize use cases
	createSiteUseCase := site.NewCreateSiteUseCase(siteRepo)
	updateSiteUseCase := site.NewUpdateSiteUseCase(siteRepo)
	deleteSiteUseCase := site.NewDeleteSiteUseCase(siteRepo)
	getSiteUseCase := site.NewGetSiteUseCase(siteRepo)
	listSitesUseCase := site.NewListSitesUseCase(siteRepo)
	executeCrawlUseCase := crawl.NewExecuteCrawlUseCase(crawlerEngine, siteRepo, crawlLogRepo)
	getLogsUseCase := crawl.NewGetLogsUseCase(crawlLogRepo)
	getLatestLogUseCase := crawl.NewGetLatestLogUseCase(crawlLogRepo)

	// Initialize handlers
	siteHandler := handler.NewSiteHandler(
		createSiteUseCase,
		updateSiteUseCase,
		deleteSiteUseCase,
		getSiteUseCase,
		listSitesUseCase,
	)
	crawlHandler := handler.NewCrawlHandler(executeCrawlUseCase)
	crawlLogHandler := handler.NewCrawlLogHandler(getLogsUseCase, getLatestLogUseCase)

	// Setup router
	r := router.SetupRouter(siteHandler, crawlHandler, crawlLogHandler)

	// Start scheduler
	ctx := context.Background()
	if err := scheduler.Start(ctx); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	// Start sync service in background (runs periodically)
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // Sync every 5 minutes
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Starting sync...")
				result, err := syncService.SyncJobs(ctx)
				if err != nil {
					log.Printf("Sync error: %v", err)
				} else {
					log.Printf("Sync completed: %d success, %d failures", result.SuccessCount, result.FailureCount)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start HTTP server
	go func() {
		addr := ":" + cfg.Port
		log.Printf("Starting HTTP server on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	scheduler.Stop()
	log.Println("Server stopped")
}


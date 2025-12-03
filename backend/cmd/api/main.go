package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authusecase "github.com/startup-job-board/backend/internal/application/usecase/auth"
	fileusecase "github.com/startup-job-board/backend/internal/application/usecase/file"
	jobusecase "github.com/startup-job-board/backend/internal/application/usecase/job"
	startupusecase "github.com/startup-job-board/backend/internal/application/usecase/startup"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/internal/infrastructure/auth"
	"github.com/startup-job-board/backend/internal/infrastructure/config"
	"github.com/startup-job-board/backend/internal/infrastructure/email"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/postgres"
	"github.com/startup-job-board/backend/internal/infrastructure/storage"
	"github.com/startup-job-board/backend/internal/presentation/http/handler"
	"github.com/startup-job-board/backend/internal/presentation/http/router"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
	"github.com/startup-job-board/backend/pkg/logger"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logger.NewLogger()

	// Initialize database
	db, err := config.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate (in production, use migrations)
	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	startupRepo := postgres.NewStartupRepository(db)
	jobRepo := postgres.NewJobRepository(db)
	memberRepo := postgres.NewStartupMemberRepository(db)
	_ = postgres.NewInvitationRepository(db) // Reserved for future use
	fileRepo := postgres.NewFileRepository(db)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWT)
	tokenGen := auth.NewTokenGenerator()
	storageService, err := storage.NewStorageService(cfg)
	if err != nil {
		logger.Warn("Failed to initialize storage: %v. File uploads will be disabled.", err)
		// Continue without storage - file uploads will fail gracefully
		storageService = nil
	}
	emailService := email.NewResendEmailService(cfg.Email)
	_ = email.NewEmailNotificationService(emailService) // Reserved for future use
	authService := service.NewAuthorizationService(memberRepo)

	// Initialize use cases
	registerUC := authusecase.NewRegisterUseCase(userRepo, logger)
	loginUC := authusecase.NewLoginUseCase(userRepo, jwtService, logger)
	refreshTokenUC := authusecase.NewRefreshTokenUseCase(userRepo, jwtService, logger)

	createStartupUC := startupusecase.NewCreateStartupUseCase(startupRepo, memberRepo, userRepo, tokenGen, logger)
	updateStartupUC := startupusecase.NewUpdateStartupUseCase(startupRepo, logger)
	getStartupUC := startupusecase.NewGetStartupUseCase(startupRepo, logger)
	listStartupsUC := startupusecase.NewListStartupsUseCase(startupRepo, logger)

	createJobUC := jobusecase.NewCreateJobUseCase(jobRepo, startupRepo, memberRepo, authService, logger)
	updateJobUC := jobusecase.NewUpdateJobUseCase(jobRepo, startupRepo, authService, logger)
	listJobsUC := jobusecase.NewListJobsUseCase(jobRepo, startupRepo, logger)
	deleteJobUC := jobusecase.NewDeleteJobUseCase(jobRepo, authService, logger)

	// Create upload use case (will handle nil storage gracefully)
	uploadFileUC := fileusecase.NewUploadFileUseCase(fileRepo, storageService, logger)

	// Initialize handlers
	validator := validator.NewValidator()
	authHandler := handler.NewAuthHandler(registerUC, loginUC, refreshTokenUC, validator)
	startupHandler := handler.NewStartupHandler(createStartupUC, updateStartupUC, getStartupUC, listStartupsUC, validator)
	jobHandler := handler.NewJobHandler(createJobUC, updateJobUC, listJobsUC, deleteJobUC, jobRepo, startupRepo, validator)
	fileHandler := handler.NewFileHandler(uploadFileUC)

	// Initialize router
	r := router.NewRouter(
		authHandler,
		startupHandler,
		jobHandler,
		fileHandler,
		jwtService,
		startupRepo,
		cfg.CORS.AllowedOrigins,
	)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Info("Server started on port %s", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func autoMigrate(db *gorm.DB) error {
	// Enable required PostgreSQL extensions
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		// Log but don't fail - extension might already exist or not be needed
		log.Printf("Warning: Could not create pgcrypto extension (might already exist): %v", err)
	}

	// GORM auto-migration - creates/updates tables based on models
	// In production, you might want to use SQL migrations instead
	return db.AutoMigrate(
		&gorm_model.User{},
		&gorm_model.Startup{},
		&gorm_model.StartupMember{},
		&gorm_model.Invitation{},
		&gorm_model.Job{},
		&gorm_model.File{},
	)
}


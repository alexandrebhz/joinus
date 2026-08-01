package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adminusecase "github.com/startup-job-board/backend/internal/application/usecase/admin"
	authusecase "github.com/startup-job-board/backend/internal/application/usecase/auth"
	billingusecase "github.com/startup-job-board/backend/internal/application/usecase/billing"
	contactusecase "github.com/startup-job-board/backend/internal/application/usecase/contact"
	fileusecase "github.com/startup-job-board/backend/internal/application/usecase/file"
	jobusecase "github.com/startup-job-board/backend/internal/application/usecase/job"
	startupusecase "github.com/startup-job-board/backend/internal/application/usecase/startup"
	teamusecase "github.com/startup-job-board/backend/internal/application/usecase/team"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/internal/infrastructure/auth"
	oauthinfra "github.com/startup-job-board/backend/internal/infrastructure/auth/oauth"
	"github.com/startup-job-board/backend/internal/infrastructure/config"
	"github.com/startup-job-board/backend/internal/infrastructure/email"
	"github.com/startup-job-board/backend/internal/infrastructure/payment"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/postgres"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/seed"
	"github.com/startup-job-board/backend/internal/infrastructure/storage"
	"github.com/startup-job-board/backend/internal/presentation/http/handler"
	"github.com/startup-job-board/backend/internal/presentation/http/router"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
	"github.com/startup-job-board/backend/pkg/logger"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.NewLogger()

	db, err := config.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	startupRepo := postgres.NewStartupRepository(db)
	jobRepo := postgres.NewJobRepository(db)
	memberRepo := postgres.NewStartupMemberRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	teamMemberRepo := postgres.NewTeamMemberRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	teamInvitationRepo := postgres.NewTeamInvitationRepository(db)
	oauthAccountRepo := postgres.NewOAuthAccountRepository(db)
	oauthLoginCodeRepo := postgres.NewOAuthLoginCodeRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	contactRepo := postgres.NewContactRepository(db)

	if err := seed.SystemRoles(context.Background(), roleRepo); err != nil {
		log.Fatalf("Failed to seed system roles: %v", err)
	}

	jwtService := auth.NewJWTService(cfg.JWT)
	tokenGen := auth.NewTokenGenerator()
	storageService, err := storage.NewStorageService(cfg)
	if err != nil {
		logger.Warn("Failed to initialize storage: %v. File uploads will be disabled.", err)
		storageService = nil
	}
	emailService := email.NewResendEmailService(cfg.Email)
	authService := service.NewAuthorizationService(userRepo, teamMemberRepo, roleRepo, startupRepo, memberRepo)
	stripeClient := payment.NewStripeClient(cfg.Stripe)

	googleRedirect := cfg.OAuth.RedirectBaseURL + "/google/callback"
	oauthRegistry := oauthinfra.NewRegistry(
		oauthinfra.NewGoogleProvider(cfg.OAuth.GoogleClientID, cfg.OAuth.GoogleClientSecret, googleRedirect),
		oauthinfra.NewStubProvider("apple"),
		oauthinfra.NewStubProvider("github"),
	)

	registerUC := authusecase.NewRegisterUseCase(userRepo, jwtService, logger)
	loginUC := authusecase.NewLoginUseCase(userRepo, jwtService, logger)
	refreshTokenUC := authusecase.NewRefreshTokenUseCase(userRepo, jwtService, logger)
	getMeUC := authusecase.NewGetMeUseCase(userRepo, teamRepo, teamMemberRepo, roleRepo, logger)
	startOAuthUC := authusecase.NewStartOAuthUseCase(oauthRegistry)
	completeOAuthUC := authusecase.NewCompleteOAuthUseCase(oauthRegistry, userRepo, oauthAccountRepo, jwtService, logger)
	issueLoginCodeUC := authusecase.NewIssueOAuthLoginCodeUseCase(oauthLoginCodeRepo)
	exchangeLoginCodeUC := authusecase.NewExchangeOAuthLoginCodeUseCase(oauthLoginCodeRepo)

	createStartupUC := startupusecase.NewCreateStartupUseCase(startupRepo, teamRepo, teamMemberRepo, roleRepo, memberRepo, userRepo, tokenGen, logger)
	updateStartupUC := startupusecase.NewUpdateStartupUseCase(startupRepo, authService, logger)
	getStartupUC := startupusecase.NewGetStartupUseCase(startupRepo, logger)
	listStartupsUC := startupusecase.NewListStartupsUseCase(startupRepo, logger)

	createJobUC := jobusecase.NewCreateJobUseCase(jobRepo, startupRepo, memberRepo, authService, logger)
	updateJobUC := jobusecase.NewUpdateJobUseCase(jobRepo, startupRepo, authService, logger)
	listJobsUC := jobusecase.NewListJobsUseCase(jobRepo, startupRepo, logger)
	deleteJobUC := jobusecase.NewDeleteJobUseCase(jobRepo, authService, logger)

	uploadFileUC := fileusecase.NewUploadFileUseCase(fileRepo, storageService, logger)
	createContactUC := contactusecase.NewCreateContactUseCase(contactRepo, logger)
	createCheckoutUC := billingusecase.NewCreateCheckoutUseCase(stripeClient, jobRepo, startupRepo, userRepo, authService, cfg.Stripe, cfg.AppURL, logger)
	handleWebhookUC := billingusecase.NewHandleWebhookUseCase(stripeClient, jobRepo, startupRepo, logger)

	createTeamUC := teamusecase.NewCreateTeamUseCase(teamRepo, teamMemberRepo, roleRepo, logger)
	listMyTeamsUC := teamusecase.NewListMyTeamsUseCase(teamRepo)
	getTeamUC := teamusecase.NewGetTeamUseCase(teamRepo, authService)
	updateTeamUC := teamusecase.NewUpdateTeamUseCase(teamRepo, authService)
	listMembersUC := teamusecase.NewListMembersUseCase(teamMemberRepo, userRepo, roleRepo, authService)
	inviteMemberUC := teamusecase.NewInviteMemberUseCase(teamInvitationRepo, teamRepo, teamMemberRepo, userRepo, roleRepo, emailService, tokenGen, authService, cfg.AppURL, logger)
	acceptInviteUC := teamusecase.NewAcceptInvitationUseCase(teamInvitationRepo, teamMemberRepo, userRepo)
	updateMemberUC := teamusecase.NewUpdateMemberUseCase(teamMemberRepo, roleRepo, authService)
	removeMemberUC := teamusecase.NewRemoveMemberUseCase(teamMemberRepo, authService)
	listRolesUC := teamusecase.NewListRolesUseCase(roleRepo, authService)
	linkStartupUC := teamusecase.NewLinkStartupUseCase(startupRepo, authService)
	unlinkStartupUC := teamusecase.NewUnlinkStartupUseCase(startupRepo, authService)
	listTeamStartupsUC := teamusecase.NewListTeamStartupsUseCase(startupRepo, authService)

	adminListUsersUC := adminusecase.NewListUsersUseCase(userRepo, authService)
	adminUpdateUserUC := adminusecase.NewUpdateUserUseCase(userRepo, authService)
	adminListTeamsUC := adminusecase.NewListTeamsUseCase(teamRepo, authService)
	adminCreateStartupUC := adminusecase.NewCreateOrphanStartupUseCase(startupRepo, tokenGen, authService, logger)
	adminLinkTeamUC := adminusecase.NewLinkStartupTeamUseCase(startupRepo, teamRepo, authService)

	v := validator.NewValidator()
	secureCookies := cfg.Environment == "production" || cfg.Environment == "prod"
	authHandler := handler.NewAuthHandler(
		registerUC, loginUC, refreshTokenUC, getMeUC,
		startOAuthUC, completeOAuthUC, issueLoginCodeUC, exchangeLoginCodeUC,
		cfg.AppURL, secureCookies, v,
	)
	startupHandler := handler.NewStartupHandler(createStartupUC, updateStartupUC, getStartupUC, listStartupsUC, v)
	jobHandler := handler.NewJobHandler(createJobUC, updateJobUC, listJobsUC, deleteJobUC, jobRepo, startupRepo, v)
	fileHandler := handler.NewFileHandler(uploadFileUC)
	contactHandler := handler.NewContactHandler(createContactUC, v)
	billingHandler := handler.NewBillingHandler(createCheckoutUC, handleWebhookUC, startupRepo, authService, v)
	teamHandler := handler.NewTeamHandler(
		createTeamUC, listMyTeamsUC, getTeamUC, updateTeamUC,
		listMembersUC, inviteMemberUC, acceptInviteUC, updateMemberUC, removeMemberUC,
		listRolesUC, linkStartupUC, unlinkStartupUC, listTeamStartupsUC, v,
	)
	adminHandler := handler.NewAdminHandler(adminListUsersUC, adminUpdateUserUC, adminListTeamsUC, adminCreateStartupUC, adminLinkTeamUC, v)

	r := router.NewRouter(router.RouterDeps{
		AuthHandler:    authHandler,
		StartupHandler: startupHandler,
		JobHandler:     jobHandler,
		FileHandler:    fileHandler,
		ContactHandler: contactHandler,
		BillingHandler: billingHandler,
		TeamHandler:    teamHandler,
		AdminHandler:   adminHandler,
		JWTService:     jwtService,
		StartupRepo:    startupRepo,
		AllowedOrigins: cfg.CORS.AllowedOrigins,
		RateLimit:      cfg.RateLimit,
		InternalKey:    cfg.InternalKey,
	})

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Info("Server started on port %s", cfg.Port)

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
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		log.Printf("Warning: Could not create pgcrypto extension (might already exist): %v", err)
	}

	return db.AutoMigrate(
		&gorm_model.User{},
		&gorm_model.Startup{},
		&gorm_model.StartupMember{},
		&gorm_model.Invitation{},
		&gorm_model.Job{},
		&gorm_model.File{},
		&gorm_model.Contact{},
		&gorm_model.Team{},
		&gorm_model.TeamMember{},
		&gorm_model.Role{},
		&gorm_model.RoleScope{},
		&gorm_model.TeamInvitation{},
		&gorm_model.OAuthAccount{},
		&gorm_model.OAuthLoginCode{},
	)
}

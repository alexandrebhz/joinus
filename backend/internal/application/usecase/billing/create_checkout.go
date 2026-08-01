package billing

import (
	"context"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/internal/infrastructure/config"
	"github.com/startup-job-board/backend/internal/infrastructure/payment"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

const (
	ProductJobBoost   = "job_boost"
	ProductStartupPro = "startup_pro"
)

// CreateCheckoutUseCase creates a Stripe Checkout Session for either a one-time
// Job Boost (€49) or a Startup Pro subscription (€99/mo).
type CreateCheckoutUseCase struct {
	stripeClient *payment.StripeClient
	jobRepo      repository.JobRepository
	startupRepo  repository.StartupRepository
	userRepo     repository.UserRepository
	authService  *service.AuthorizationService
	stripeCfg    config.StripeConfig
	appURL       string
	logger       logger.Logger
}

func NewCreateCheckoutUseCase(
	stripeClient *payment.StripeClient,
	jobRepo repository.JobRepository,
	startupRepo repository.StartupRepository,
	userRepo repository.UserRepository,
	authService *service.AuthorizationService,
	stripeCfg config.StripeConfig,
	appURL string,
	logger logger.Logger,
) *CreateCheckoutUseCase {
	return &CreateCheckoutUseCase{
		stripeClient: stripeClient,
		jobRepo:      jobRepo,
		startupRepo:  startupRepo,
		userRepo:     userRepo,
		authService:  authService,
		stripeCfg:    stripeCfg,
		appURL:       appURL,
		logger:       logger,
	}
}

func (uc *CreateCheckoutUseCase) Execute(ctx context.Context, input dto.CreateCheckoutInput, userID string) (*dto.CheckoutOutput, error) {
	switch input.Product {
	case ProductJobBoost:
		return uc.checkoutJobBoost(ctx, input, userID)
	case ProductStartupPro:
		return uc.checkoutStartupPro(ctx, input, userID)
	default:
		return nil, errors.NewBadRequestError("unsupported product")
	}
}

func (uc *CreateCheckoutUseCase) checkoutJobBoost(ctx context.Context, input dto.CreateCheckoutInput, userID string) (*dto.CheckoutOutput, error) {
	if input.JobID == nil || *input.JobID == "" {
		return nil, errors.NewBadRequestError("job_id is required for job_boost")
	}

	job, err := uc.jobRepo.FindByID(ctx, *input.JobID)
	if err != nil {
		return nil, errors.NewNotFoundError("job")
	}

	canManage, err := uc.authService.CanManageJobs(ctx, userID, job.StartupID)
	if err != nil || !canManage {
		return nil, errors.NewForbiddenError("you don't have permission to boost this job")
	}

	if uc.stripeCfg.PriceJobBoost == "" {
		return nil, errors.NewBadRequestError("job boost is not configured")
	}

	user, _ := uc.userRepo.FindByID(ctx, userID)
	customerEmail := ""
	if user != nil {
		customerEmail = user.Email
	}

	result, err := uc.stripeClient.CreateCheckoutSession(payment.CheckoutSessionInput{
		Mode:          "payment",
		PriceID:       uc.stripeCfg.PriceJobBoost,
		SuccessURL:    uc.appURL + "/billing/success?session_id={CHECKOUT_SESSION_ID}",
		CancelURL:     uc.appURL + "/billing/cancel",
		CustomerEmail: customerEmail,
		Metadata: map[string]string{
			"product":    ProductJobBoost,
			"job_id":     job.ID,
			"startup_id": job.StartupID,
			"user_id":    userID,
		},
	})
	if err != nil {
		uc.logger.Error("Failed to create job boost checkout session: %v", err)
		return nil, errors.NewBadRequestError("failed to create checkout session")
	}

	return &dto.CheckoutOutput{URL: result.URL, SessionID: result.ID}, nil
}

func (uc *CreateCheckoutUseCase) checkoutStartupPro(ctx context.Context, input dto.CreateCheckoutInput, userID string) (*dto.CheckoutOutput, error) {
	if input.StartupID == nil || *input.StartupID == "" {
		return nil, errors.NewBadRequestError("startup_id is required for startup_pro")
	}

	startup, err := uc.startupRepo.FindByID(ctx, *input.StartupID)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

	canManage, err := uc.authService.CanManageStartup(ctx, userID, startup.ID)
	if err != nil || !canManage {
		return nil, errors.NewForbiddenError("you don't have permission to manage billing for this startup")
	}

	if startup.Plan == string(entity.StartupPlanPro) {
		return nil, errors.NewBadRequestError("startup is already on the Pro plan")
	}

	if uc.stripeCfg.PriceStartupPro == "" {
		return nil, errors.NewBadRequestError("startup pro is not configured")
	}

	sessionInput := payment.CheckoutSessionInput{
		Mode:       "subscription",
		PriceID:    uc.stripeCfg.PriceStartupPro,
		SuccessURL: uc.appURL + "/billing/success?session_id={CHECKOUT_SESSION_ID}",
		CancelURL:  uc.appURL + "/billing/cancel",
		Metadata: map[string]string{
			"product":    ProductStartupPro,
			"startup_id": startup.ID,
			"user_id":    userID,
		},
		SubscriptionMetadata: map[string]string{
			"product":    ProductStartupPro,
			"startup_id": startup.ID,
		},
	}

	if startup.StripeCustomerID != nil && *startup.StripeCustomerID != "" {
		sessionInput.CustomerID = *startup.StripeCustomerID
	} else {
		user, _ := uc.userRepo.FindByID(ctx, userID)
		if user != nil {
			sessionInput.CustomerEmail = user.Email
		}
	}

	result, err := uc.stripeClient.CreateCheckoutSession(sessionInput)
	if err != nil {
		uc.logger.Error("Failed to create startup pro checkout session: %v", err)
		return nil, errors.NewBadRequestError("failed to create checkout session")
	}

	return &dto.CheckoutOutput{URL: result.URL, SessionID: result.ID}, nil
}

package billing

import (
	"context"
	"encoding/json"
	"time"

	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/payment"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
	"github.com/stripe/stripe-go/v82"
)

const (
	jobBoostDuration   = 30 * 24 * time.Hour
	startupProDuration = 30 * 24 * time.Hour
)

// HandleWebhookUseCase applies Stripe billing events to the domain:
//   - checkout.session.completed: provisions the purchased product (job boost or startup pro)
//   - customer.subscription.updated/deleted: keeps the Startup Pro plan status in sync
//   - invoice.payment_succeeded: extends the Startup Pro plan on subscription renewal
type HandleWebhookUseCase struct {
	stripeClient *payment.StripeClient
	jobRepo      repository.JobRepository
	startupRepo  repository.StartupRepository
	logger       logger.Logger
}

func NewHandleWebhookUseCase(
	stripeClient *payment.StripeClient,
	jobRepo repository.JobRepository,
	startupRepo repository.StartupRepository,
	logger logger.Logger,
) *HandleWebhookUseCase {
	return &HandleWebhookUseCase{
		stripeClient: stripeClient,
		jobRepo:      jobRepo,
		startupRepo:  startupRepo,
		logger:       logger,
	}
}

func (uc *HandleWebhookUseCase) Execute(ctx context.Context, payload []byte, signatureHeader string) error {
	event, err := uc.stripeClient.ConstructWebhookEvent(payload, signatureHeader)
	if err != nil {
		return errors.NewBadRequestError("invalid webhook signature")
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		return uc.handleCheckoutCompleted(ctx, event)
	case stripe.EventTypeCustomerSubscriptionUpdated:
		return uc.handleSubscriptionUpdated(ctx, event)
	case stripe.EventTypeCustomerSubscriptionDeleted:
		return uc.handleSubscriptionDeleted(ctx, event)
	case stripe.EventTypeInvoicePaymentSucceeded:
		return uc.handleInvoicePaymentSucceeded(ctx, event)
	default:
		// Unhandled event types are acknowledged without action.
		return nil
	}
}

func (uc *HandleWebhookUseCase) handleCheckoutCompleted(ctx context.Context, event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		uc.logger.Error("Failed to parse checkout session: %v", err)
		return err
	}

	switch session.Metadata["product"] {
	case ProductJobBoost:
		return uc.provisionJobBoost(ctx, session)
	case ProductStartupPro:
		return uc.provisionStartupPro(ctx, session)
	default:
		uc.logger.Warn("Checkout session completed with unknown product metadata: %v", session.Metadata)
		return nil
	}
}

func (uc *HandleWebhookUseCase) provisionJobBoost(ctx context.Context, session stripe.CheckoutSession) error {
	jobID := session.Metadata["job_id"]
	if jobID == "" {
		return nil
	}

	job, err := uc.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		uc.logger.Error("Job boost webhook: job %s not found: %v", jobID, err)
		return err
	}

	boostedUntil := time.Now().Add(jobBoostDuration)
	job.BoostedUntil = &boostedUntil
	job.UpdatedAt = time.Now()

	if err := uc.jobRepo.Update(ctx, job); err != nil {
		uc.logger.Error("Failed to update job boost for %s: %v", jobID, err)
		return err
	}

	uc.logger.Info("Job %s boosted until %s", jobID, boostedUntil.Format(time.RFC3339))
	return nil
}

func (uc *HandleWebhookUseCase) provisionStartupPro(ctx context.Context, session stripe.CheckoutSession) error {
	startupID := session.Metadata["startup_id"]
	if startupID == "" {
		return nil
	}

	startup, err := uc.startupRepo.FindByID(ctx, startupID)
	if err != nil {
		uc.logger.Error("Startup pro webhook: startup %s not found: %v", startupID, err)
		return err
	}

	planExpiresAt := time.Now().Add(startupProDuration)
	startup.Plan = string(entity.StartupPlanPro)
	startup.PlanExpiresAt = &planExpiresAt
	if session.Customer != nil && session.Customer.ID != "" {
		startup.StripeCustomerID = &session.Customer.ID
	}
	if session.Subscription != nil && session.Subscription.ID != "" {
		startup.StripeSubscriptionID = &session.Subscription.ID
	}
	startup.UpdatedAt = time.Now()

	if err := uc.startupRepo.Update(ctx, startup); err != nil {
		uc.logger.Error("Failed to activate startup pro for %s: %v", startupID, err)
		return err
	}

	uc.logger.Info("Startup %s upgraded to Pro until %s", startupID, planExpiresAt.Format(time.RFC3339))
	return nil
}

func (uc *HandleWebhookUseCase) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		uc.logger.Error("Failed to parse subscription: %v", err)
		return err
	}

	startup, err := uc.startupRepo.FindByStripeSubscriptionID(ctx, sub.ID)
	if err != nil {
		// Subscription not tied to a startup (or not yet provisioned) - nothing to do.
		return nil
	}

	switch sub.Status {
	case stripe.SubscriptionStatusActive, stripe.SubscriptionStatusTrialing:
		planExpiresAt := time.Now().Add(startupProDuration)
		startup.Plan = string(entity.StartupPlanPro)
		startup.PlanExpiresAt = &planExpiresAt
	case stripe.SubscriptionStatusCanceled, stripe.SubscriptionStatusUnpaid, stripe.SubscriptionStatusIncompleteExpired:
		startup.Plan = string(entity.StartupPlanFree)
		startup.PlanExpiresAt = nil
	}
	startup.UpdatedAt = time.Now()

	if err := uc.startupRepo.Update(ctx, startup); err != nil {
		uc.logger.Error("Failed to sync subscription update for startup %s: %v", startup.ID, err)
		return err
	}

	return nil
}

func (uc *HandleWebhookUseCase) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		uc.logger.Error("Failed to parse subscription: %v", err)
		return err
	}

	startup, err := uc.startupRepo.FindByStripeSubscriptionID(ctx, sub.ID)
	if err != nil {
		return nil
	}

	startup.Plan = string(entity.StartupPlanFree)
	startup.PlanExpiresAt = nil
	startup.UpdatedAt = time.Now()

	if err := uc.startupRepo.Update(ctx, startup); err != nil {
		uc.logger.Error("Failed to downgrade startup %s after subscription deletion: %v", startup.ID, err)
		return err
	}

	uc.logger.Info("Startup %s downgraded to Free (subscription canceled)", startup.ID)
	return nil
}

func (uc *HandleWebhookUseCase) handleInvoicePaymentSucceeded(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		uc.logger.Error("Failed to parse invoice: %v", err)
		return err
	}

	subscriptionID := invoiceSubscriptionID(&invoice)
	if subscriptionID == "" {
		return nil
	}

	startup, err := uc.startupRepo.FindByStripeSubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil
	}

	// Renewal payment succeeded: extend the plan for another billing cycle.
	planExpiresAt := time.Now().Add(startupProDuration)
	startup.Plan = string(entity.StartupPlanPro)
	startup.PlanExpiresAt = &planExpiresAt
	startup.UpdatedAt = time.Now()

	if err := uc.startupRepo.Update(ctx, startup); err != nil {
		uc.logger.Error("Failed to extend startup pro plan for %s: %v", startup.ID, err)
		return err
	}

	return nil
}

// invoiceSubscriptionID extracts the related subscription ID from an invoice,
// accounting for the API's move of this reference under Parent.SubscriptionDetails.
func invoiceSubscriptionID(invoice *stripe.Invoice) string {
	if invoice.Parent != nil && invoice.Parent.SubscriptionDetails != nil && invoice.Parent.SubscriptionDetails.Subscription != nil {
		return invoice.Parent.SubscriptionDetails.Subscription.ID
	}
	return ""
}

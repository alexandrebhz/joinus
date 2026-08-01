package payment

import (
	"github.com/startup-job-board/backend/internal/infrastructure/config"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
)

// StripeClient wraps the Stripe SDK for Checkout session creation and webhook
// signature verification.
//
// Plan mapping:
//   - Job Boost (one-time, €49)          -> "payment" mode checkout
//   - Startup Pro (subscription, €99/mo) -> "subscription" mode checkout
type StripeClient struct {
	webhookSecret string
}

func NewStripeClient(cfg config.StripeConfig) *StripeClient {
	stripe.Key = cfg.SecretKey
	return &StripeClient{webhookSecret: cfg.WebhookSecret}
}

type CheckoutSessionInput struct {
	Mode                 string // "payment" | "subscription"
	PriceID              string
	SuccessURL           string
	CancelURL            string
	CustomerID           string // existing Stripe customer ID, takes precedence over CustomerEmail
	CustomerEmail        string
	Metadata             map[string]string // attached to the Checkout Session (and, for one-time payments, the resulting object)
	SubscriptionMetadata map[string]string // attached to the created Subscription (subscription mode only)
}

type CheckoutSessionOutput struct {
	ID  string
	URL string
}

func (s *StripeClient) CreateCheckoutSession(input CheckoutSessionInput) (*CheckoutSessionOutput, error) {
	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(input.Mode),
		SuccessURL: stripe.String(input.SuccessURL),
		CancelURL:  stripe.String(input.CancelURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(input.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: input.Metadata,
	}

	if input.CustomerID != "" {
		params.Customer = stripe.String(input.CustomerID)
	} else if input.CustomerEmail != "" {
		params.CustomerEmail = stripe.String(input.CustomerEmail)
	}

	if input.Mode == string(stripe.CheckoutSessionModeSubscription) && len(input.SubscriptionMetadata) > 0 {
		params.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: input.SubscriptionMetadata,
		}
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return &CheckoutSessionOutput{ID: sess.ID, URL: sess.URL}, nil
}

// ConstructWebhookEvent verifies the Stripe-Signature header and decodes the event payload.
func (s *StripeClient) ConstructWebhookEvent(payload []byte, signatureHeader string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signatureHeader, s.webhookSecret)
}

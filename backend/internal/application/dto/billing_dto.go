package dto

// CreateCheckoutInput is the request body for POST /billing/checkout.
//
// Product mapping:
//   - "job_boost"   -> Job Boost, one-time €49, requires JobID, boosts the job for 30 days
//   - "startup_pro" -> Startup Pro, €99/mo subscription, requires StartupID
type CreateCheckoutInput struct {
	Product   string  `json:"product" validate:"required,oneof=job_boost startup_pro"`
	JobID     *string `json:"job_id"`
	StartupID *string `json:"startup_id"`
}

type CheckoutOutput struct {
	URL       string `json:"url"`
	SessionID string `json:"session_id"`
}

type StartupBillingStatus struct {
	StartupID     string  `json:"startup_id"`
	StartupName   string  `json:"startup_name"`
	Plan          string  `json:"plan"`
	PlanExpiresAt *string `json:"plan_expires_at"`
}

type JobBoostStatus struct {
	JobID        string  `json:"job_id"`
	JobTitle     string  `json:"job_title"`
	Boosted      bool    `json:"boosted"`
	BoostedUntil *string `json:"boosted_until"`
}

type BillingStatusOutput struct {
	Startups []StartupBillingStatus `json:"startups"`
}

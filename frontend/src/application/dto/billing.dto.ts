export type BillingProduct = 'job_boost' | 'startup_pro'

export interface CreateCheckoutRequest {
  product: BillingProduct
  job_id?: string
  startup_id?: string
}

export interface CheckoutResponse {
  url: string
  sessionId: string
}

export interface StartupBillingStatus {
  startupId: string
  startupName: string
  plan: 'free' | 'pro'
  planExpiresAt?: string
}

export interface BillingStatusResponse {
  startups: StartupBillingStatus[]
}

package entity

import "time"

type OAuthProvider string

const (
	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderApple  OAuthProvider = "apple"
	OAuthProviderGitHub OAuthProvider = "github"
)

type OAuthAccount struct {
	ID             string
	UserID         string
	Provider       OAuthProvider
	ProviderUserID string
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

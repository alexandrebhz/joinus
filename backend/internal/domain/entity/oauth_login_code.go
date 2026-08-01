package entity

import "time"

// OAuthLoginCode is a single-use bridge between the OAuth callback and the SPA.
// The browser only ever sees the opaque Code — never JWTs in the URL.
type OAuthLoginCode struct {
	ID           string
	CodeHash     string
	AccessToken  string
	RefreshToken string
	UserID       string
	ExpiresAt    time.Time
	UsedAt       *time.Time
	CreatedAt    time.Time
}

func (c *OAuthLoginCode) IsValid(now time.Time) bool {
	return c.UsedAt == nil && now.Before(c.ExpiresAt)
}

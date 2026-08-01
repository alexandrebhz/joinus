package port

import "context"

// ExternalProfile is the normalized identity returned by an OAuth provider.
type ExternalProfile struct {
	ProviderUserID string
	Email          string
	Name           string
}

// OAuthProvider abstracts Google / Apple / GitHub (etc.) auth-code flows.
type OAuthProvider interface {
	Name() string
	AuthURL(state string) string
	Exchange(ctx context.Context, code string) (*ExternalProfile, error)
}

// OAuthRegistry looks up a provider by name.
type OAuthRegistry interface {
	Get(name string) (OAuthProvider, bool)
}

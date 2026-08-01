package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/pkg/errors"
)

type Registry struct {
	mu        sync.RWMutex
	providers map[string]port.OAuthProvider
}

func NewRegistry(providers ...port.OAuthProvider) *Registry {
	r := &Registry{providers: map[string]port.OAuthProvider{}}
	for _, p := range providers {
		if p != nil {
			r.providers[p.Name()] = p
		}
	}
	return r
}

func (r *Registry) Get(name string) (port.OAuthProvider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	return p, ok
}

type GoogleProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	httpClient   *http.Client
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) *GoogleProvider {
	if clientID == "" || clientSecret == "" {
		return nil
	}
	return &GoogleProvider{
		clientID: clientID, clientSecret: clientSecret, redirectURL: redirectURL,
		httpClient: http.DefaultClient,
	}
}

func (p *GoogleProvider) Name() string { return "google" }

func (p *GoogleProvider) AuthURL(state string) string {
	q := url.Values{}
	q.Set("client_id", p.clientID)
	q.Set("redirect_uri", p.redirectURL)
	q.Set("response_type", "code")
	q.Set("scope", "openid email profile")
	q.Set("state", state)
	q.Set("access_type", "online")
	q.Set("prompt", "select_account")
	return "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode()
}

func (p *GoogleProvider) Exchange(ctx context.Context, code string) (*port.ExternalProfile, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", p.clientID)
	form.Set("client_secret", p.clientSecret)
	form.Set("redirect_uri", p.redirectURL)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google token exchange failed: %s", string(body))
	}
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	ureq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return nil, err
	}
	ureq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	uresp, err := p.httpClient.Do(ureq)
	if err != nil {
		return nil, err
	}
	defer uresp.Body.Close()
	ubody, _ := io.ReadAll(uresp.Body)
	if uresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo failed: %s", string(ubody))
	}
	var info struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(ubody, &info); err != nil {
		return nil, err
	}
	return &port.ExternalProfile{
		ProviderUserID: info.Sub,
		Email:          info.Email,
		Name:           info.Name,
	}, nil
}

// StubProvider returns a clear error for unimplemented providers (Apple, GitHub).
type StubProvider struct {
	name string
}

func NewStubProvider(name string) *StubProvider {
	return &StubProvider{name: name}
}

func (p *StubProvider) Name() string { return p.name }

func (p *StubProvider) AuthURL(state string) string { return "" }

func (p *StubProvider) Exchange(ctx context.Context, code string) (*port.ExternalProfile, error) {
	return nil, errors.NewBadRequestError(p.name + " oauth is not implemented yet")
}

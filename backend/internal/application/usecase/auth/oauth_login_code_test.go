package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/startup-job-board/backend/internal/application/dto"
	authusecase "github.com/startup-job-board/backend/internal/application/usecase/auth"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type memCodeRepo struct {
	byHash map[string]*entity.OAuthLoginCode
}

func (r *memCodeRepo) Create(ctx context.Context, code *entity.OAuthLoginCode) error {
	r.byHash[code.CodeHash] = code
	return nil
}

func (r *memCodeRepo) FindByCodeHash(ctx context.Context, codeHash string) (*entity.OAuthLoginCode, error) {
	c, ok := r.byHash[codeHash]
	if !ok {
		return nil, nil
	}
	cp := *c
	return &cp, nil
}

func (r *memCodeRepo) MarkUsed(ctx context.Context, id string, usedAt time.Time) error {
	for _, c := range r.byHash {
		if c.ID == id {
			c.UsedAt = &usedAt
			return nil
		}
	}
	return errors.New("not found")
}

func (r *memCodeRepo) DeleteExpired(ctx context.Context) error { return nil }

type memUserRepo struct {
	users map[string]*entity.User
}

func (r *memUserRepo) Create(ctx context.Context, user *entity.User) error { return nil }
func (r *memUserRepo) Update(ctx context.Context, user *entity.User) error { return nil }
func (r *memUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (r *memUserRepo) List(ctx context.Context, page, pageSize int, search string) ([]*entity.User, int64, error) {
	return nil, 0, nil
}
func (r *memUserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	cp := *u
	return &cp, nil
}

type mockJWT struct{}

func (m *mockJWT) GenerateAccessToken(userID, role string) (string, error) {
	return "access-" + userID, nil
}
func (m *mockJWT) GenerateRefreshToken(userID string, version int) (string, error) {
	return "refresh-" + userID, nil
}
func (m *mockJWT) ValidateToken(token string) (string, string, error) { return "", "", nil }
func (m *mockJWT) ValidateRefreshToken(token string) (string, int, error) {
	return "", 0, nil
}

func TestOAuthLoginCodeIsSingleUseAndNotInTokens(t *testing.T) {
	repo := &memCodeRepo{byHash: map[string]*entity.OAuthLoginCode{}}
	userRepo := &memUserRepo{users: map[string]*entity.User{
		"user-1": {
			ID: "user-1", Email: "a@b.com", Name: "A",
			Role: entity.UserRoleCandidate, Status: entity.UserStatusActive,
		},
	}}
	jwt := &mockJWT{}
	issue := authusecase.NewIssueOAuthLoginCodeUseCase(repo)
	exchange := authusecase.NewExchangeOAuthLoginCodeUseCase(repo, userRepo, jwt)

	authOut := &dto.AuthOutput{
		AccessToken:  "access-jwt",
		RefreshToken: "refresh-jwt",
		User:         dto.UserOutput{ID: "user-1", Email: "a@b.com", Name: "A", Role: "candidate"},
	}

	code, err := issue.Execute(context.Background(), authOut)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if code == "" {
		t.Fatal("expected opaque code")
	}
	if code == authOut.AccessToken || code == authOut.RefreshToken {
		t.Fatal("login code must not be a JWT")
	}

	stored := findStoredCode(repo)
	if stored == nil {
		t.Fatal("expected stored code")
	}
	if stored.AccessToken != "" || stored.RefreshToken != "" {
		t.Fatal("stored code must not contain JWTs")
	}

	first, err := exchange.Execute(context.Background(), dto.ExchangeOAuthCodeInput{Code: code})
	if err != nil {
		t.Fatalf("first exchange: %v", err)
	}
	if first.AccessToken != "access-user-1" || first.RefreshToken != "refresh-user-1" {
		t.Fatalf("unexpected tokens: %+v", first)
	}

	_, err = exchange.Execute(context.Background(), dto.ExchangeOAuthCodeInput{Code: code})
	if err == nil {
		t.Fatal("second exchange must fail (single use)")
	}
}

func TestOAuthLoginCodeRejectsUnknown(t *testing.T) {
	repo := &memCodeRepo{byHash: map[string]*entity.OAuthLoginCode{}}
	userRepo := &memUserRepo{users: map[string]*entity.User{}}
	exchange := authusecase.NewExchangeOAuthLoginCodeUseCase(repo, userRepo, &mockJWT{})
	_, err := exchange.Execute(context.Background(), dto.ExchangeOAuthCodeInput{Code: "nope"})
	if err == nil {
		t.Fatal("expected error for unknown code")
	}
}

func findStoredCode(repo *memCodeRepo) *entity.OAuthLoginCode {
	for _, c := range repo.byHash {
		return c
	}
	return nil
}

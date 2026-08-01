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
	// return a copy so mutations don't surprise us
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

func TestOAuthLoginCodeIsSingleUseAndNotInTokens(t *testing.T) {
	repo := &memCodeRepo{byHash: map[string]*entity.OAuthLoginCode{}}
	issue := authusecase.NewIssueOAuthLoginCodeUseCase(repo)
	exchange := authusecase.NewExchangeOAuthLoginCodeUseCase(repo)

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

	first, err := exchange.Execute(context.Background(), dto.ExchangeOAuthCodeInput{Code: code})
	if err != nil {
		t.Fatalf("first exchange: %v", err)
	}
	if first.AccessToken != "access-jwt" || first.RefreshToken != "refresh-jwt" {
		t.Fatalf("unexpected tokens: %+v", first)
	}

	_, err = exchange.Execute(context.Background(), dto.ExchangeOAuthCodeInput{Code: code})
	if err == nil {
		t.Fatal("second exchange must fail (single use)")
	}
}

func TestOAuthLoginCodeRejectsUnknown(t *testing.T) {
	repo := &memCodeRepo{byHash: map[string]*entity.OAuthLoginCode{}}
	exchange := authusecase.NewExchangeOAuthLoginCodeUseCase(repo)
	_, err := exchange.Execute(context.Background(), dto.ExchangeOAuthCodeInput{Code: "nope"})
	if err == nil {
		t.Fatal("expected error for unknown code")
	}
}

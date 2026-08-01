package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
)

const oauthLoginCodeTTL = 60 * time.Second

// IssueOAuthLoginCodeUseCase stores JWTs server-side and returns an opaque one-time code
// safe to put in a redirect URL.
type IssueOAuthLoginCodeUseCase struct {
	codeRepo repository.OAuthLoginCodeRepository
}

func NewIssueOAuthLoginCodeUseCase(codeRepo repository.OAuthLoginCodeRepository) *IssueOAuthLoginCodeUseCase {
	return &IssueOAuthLoginCodeUseCase{codeRepo: codeRepo}
}

func (uc *IssueOAuthLoginCodeUseCase) Execute(ctx context.Context, auth *dto.AuthOutput) (plainCode string, err error) {
	plain, err := randomOAuthCode()
	if err != nil {
		return "", errors.ErrInternalError
	}
	now := time.Now()
	record := &entity.OAuthLoginCode{
		ID:           uuid.New().String(),
		CodeHash:     hashOAuthCode(plain),
		AccessToken:  auth.AccessToken,
		RefreshToken: auth.RefreshToken,
		UserID:       auth.User.ID,
		ExpiresAt:    now.Add(oauthLoginCodeTTL),
		CreatedAt:    now,
	}
	if err := uc.codeRepo.Create(ctx, record); err != nil {
		return "", err
	}
	_ = uc.codeRepo.DeleteExpired(ctx)
	return plain, nil
}

// ExchangeOAuthLoginCodeUseCase consumes a one-time code and returns the JWT pair once.
type ExchangeOAuthLoginCodeUseCase struct {
	codeRepo repository.OAuthLoginCodeRepository
}

func NewExchangeOAuthLoginCodeUseCase(codeRepo repository.OAuthLoginCodeRepository) *ExchangeOAuthLoginCodeUseCase {
	return &ExchangeOAuthLoginCodeUseCase{codeRepo: codeRepo}
}

func (uc *ExchangeOAuthLoginCodeUseCase) Execute(ctx context.Context, input dto.ExchangeOAuthCodeInput) (*dto.AuthOutput, error) {
	if input.Code == "" {
		return nil, errors.NewBadRequestError("invalid or expired code")
	}
	record, err := uc.codeRepo.FindByCodeHash(ctx, hashOAuthCode(input.Code))
	if err != nil {
		return nil, err
	}
	if record == nil || !record.IsValid(time.Now()) {
		return nil, errors.NewUnauthorizedError("invalid or expired code")
	}

	now := time.Now()
	if err := uc.codeRepo.MarkUsed(ctx, record.ID, now); err != nil {
		return nil, errors.NewUnauthorizedError("invalid or expired code")
	}

	return &dto.AuthOutput{
		AccessToken:  record.AccessToken,
		RefreshToken: record.RefreshToken,
		User:         dto.UserOutput{ID: record.UserID},
	}, nil
}

func randomOAuthCode() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashOAuthCode(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

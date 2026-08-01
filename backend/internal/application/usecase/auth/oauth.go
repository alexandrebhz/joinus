package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type StartOAuthUseCase struct {
	registry port.OAuthRegistry
}

func NewStartOAuthUseCase(registry port.OAuthRegistry) *StartOAuthUseCase {
	return &StartOAuthUseCase{registry: registry}
}

func (uc *StartOAuthUseCase) Execute(providerName string) (authURL, state string, err error) {
	provider, ok := uc.registry.Get(providerName)
	if !ok {
		return "", "", errors.NewBadRequestError("unsupported oauth provider")
	}
	state, err = randomState()
	if err != nil {
		return "", "", errors.ErrInternalError
	}
	return provider.AuthURL(state), state, nil
}

type CompleteOAuthUseCase struct {
	registry   port.OAuthRegistry
	userRepo   repository.UserRepository
	oauthRepo  repository.OAuthAccountRepository
	jwtService port.JWTService
	logger     logger.Logger
}

func NewCompleteOAuthUseCase(
	registry port.OAuthRegistry,
	userRepo repository.UserRepository,
	oauthRepo repository.OAuthAccountRepository,
	jwtService port.JWTService,
	logger logger.Logger,
) *CompleteOAuthUseCase {
	return &CompleteOAuthUseCase{
		registry: registry, userRepo: userRepo, oauthRepo: oauthRepo,
		jwtService: jwtService, logger: logger,
	}
}

func (uc *CompleteOAuthUseCase) Execute(ctx context.Context, providerName, code string) (*dto.AuthOutput, error) {
	provider, ok := uc.registry.Get(providerName)
	if !ok {
		return nil, errors.NewBadRequestError("unsupported oauth provider")
	}
	profile, err := provider.Exchange(ctx, code)
	if err != nil {
		uc.logger.Warn("oauth exchange failed for %s: %v", providerName, err)
		return nil, errors.NewUnauthorizedError("oauth authentication failed")
	}
	if profile.Email == "" || profile.ProviderUserID == "" {
		return nil, errors.NewBadRequestError("oauth provider returned incomplete profile")
	}

	prov := entity.OAuthProvider(providerName)
	account, err := uc.oauthRepo.FindByProviderAndUserID(ctx, prov, profile.ProviderUserID)
	if err != nil {
		return nil, err
	}

	var user *entity.User
	if account != nil {
		user, err = uc.userRepo.FindByID(ctx, account.UserID)
		if err != nil {
			return nil, errors.NewUnauthorizedError("oauth authentication failed")
		}
	} else {
		user, _ = uc.userRepo.FindByEmail(ctx, profile.Email)
		if user == nil {
			// Create user with random unusable password (OAuth-only).
			hash, _ := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
			name := profile.Name
			if name == "" {
				name = profile.Email
			}
			now := time.Now()
			user = &entity.User{
				ID: uuid.New().String(), Email: profile.Email, Password: string(hash),
				Name: name, Role: entity.UserRoleCandidate, Status: entity.UserStatusActive,
				CreatedAt: now, UpdatedAt: now,
			}
			if err := uc.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		}
		now := time.Now()
		if err := uc.oauthRepo.Create(ctx, &entity.OAuthAccount{
			ID: uuid.New().String(), UserID: user.ID, Provider: prov,
			ProviderUserID: profile.ProviderUserID, Email: profile.Email,
			CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			return nil, err
		}
	}

	if user.Status != entity.UserStatusActive {
		return nil, errors.NewUnauthorizedError("account is not active")
	}

	access, err := uc.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}
	refresh, err := uc.jwtService.GenerateRefreshToken(user.ID, user.TokenVersion)
	if err != nil {
		return nil, err
	}
	return &dto.AuthOutput{
		AccessToken: access, RefreshToken: refresh,
		User: dto.UserOutput{ID: user.ID, Email: user.Email, Name: user.Name, Role: string(user.Role)},
	}, nil
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

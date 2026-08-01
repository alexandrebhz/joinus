package auth

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/pkg/logger"
)

type RefreshTokenUseCase struct {
	userRepo   repository.UserRepository
	jwtService port.JWTService
	logger     logger.Logger
}

func NewRefreshTokenUseCase(
	userRepo repository.UserRepository,
	jwtService port.JWTService,
	logger logger.Logger,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (*dto.AuthOutput, error) {
	userID, version, err := uc.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	if user.Status != "active" {
		return nil, errors.NewForbiddenError("account is not active")
	}

	if user.TokenVersion != version {
		return nil, errors.ErrUnauthorized
	}

	user.TokenVersion++
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID, user.TokenVersion)
	if err != nil {
		return nil, err
	}

	return &dto.AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: dto.UserOutput{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  string(user.Role),
		},
	}, nil
}




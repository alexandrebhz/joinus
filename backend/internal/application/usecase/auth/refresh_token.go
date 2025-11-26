package auth

import (
	"context"

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
	// Validate refresh token
	userID, err := uc.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	// Get user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	if user.Status != "active" {
		return nil, errors.NewForbiddenError("account is not active")
	}

	// Generate new tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
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




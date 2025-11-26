package auth

import (
	"context"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	userRepo     repository.UserRepository
	jwtService   port.JWTService
	logger       logger.Logger
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	jwtService port.JWTService,
	logger logger.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input dto.LoginInput) (*dto.AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	if user.Status != "active" {
		return nil, errors.NewForbiddenError("account is not active")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.ErrUnauthorized
	}

	// Generate tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserOutput{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  string(user.Role),
		},
	}, nil
}




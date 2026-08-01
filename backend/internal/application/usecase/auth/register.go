package auth

import (
	"context"
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

type RegisterUseCase struct {
	userRepo   repository.UserRepository
	jwtService port.JWTService
	logger     logger.Logger
}

func NewRegisterUseCase(
	userRepo repository.UserRepository,
	jwtService port.JWTService,
	logger logger.Logger,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input dto.RegisterInput) (*dto.AuthOutput, error) {
	existing, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.NewBadRequestError("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternalError
	}

	user := &entity.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		Password:  string(hashedPassword),
		Name:      input.Name,
		Role:      entity.UserRoleCandidate,
		Status:    entity.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	access, err := uc.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}
	refresh, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthOutput{
		AccessToken:  access,
		RefreshToken: refresh,
		User: dto.UserOutput{
			ID: user.ID, Email: user.Email, Name: user.Name, Role: string(user.Role),
		},
	}, nil
}

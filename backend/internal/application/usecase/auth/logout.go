package auth

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type LogoutUseCase struct {
	userRepo repository.UserRepository
	logger   logger.Logger
}

func NewLogoutUseCase(userRepo repository.UserRepository, logger logger.Logger) *LogoutUseCase {
	return &LogoutUseCase{userRepo: userRepo, logger: logger}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, userID string) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.ErrUnauthorized
	}

	user.TokenVersion++
	user.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, user)
}

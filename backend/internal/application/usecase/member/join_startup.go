package member

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
	"github.com/startup-job-board/backend/pkg/utils"
)

type JoinStartupUseCase struct {
	startupRepo repository.StartupRepository
	memberRepo  repository.StartupMemberRepository
	userRepo    repository.UserRepository
	logger      logger.Logger
}

func NewJoinStartupUseCase(
	startupRepo repository.StartupRepository,
	memberRepo repository.StartupMemberRepository,
	userRepo repository.UserRepository,
	logger logger.Logger,
) *JoinStartupUseCase {
	return &JoinStartupUseCase{
		startupRepo: startupRepo,
		memberRepo:  memberRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (uc *JoinStartupUseCase) Execute(ctx context.Context, startupSlug string, joinCode string, userID string) (*dto.MemberOutput, error) {
	// Find startup by slug
	startup, err := uc.startupRepo.FindBySlug(ctx, startupSlug)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

	// Check if public join is allowed
	if !startup.AllowPublicJoin {
		return nil, errors.NewForbiddenError("this startup does not allow public joining")
	}

	// Validate join code if required
	if startup.JoinCode != nil && *startup.JoinCode != "" {
		// Hash the provided code and compare
		providedHash := utils.HashString(joinCode)
		expectedHash := utils.HashString(*startup.JoinCode)
		if providedHash != expectedHash {
			return nil, errors.NewForbiddenError("invalid join code")
		}
	}

	// Check if already a member
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}

	existing, _ := uc.memberRepo.FindByUserAndStartup(ctx, userID, startup.ID)
	if existing != nil && existing.Status == entity.MemberStatusActive {
		return nil, errors.NewBadRequestError("you are already a member of this startup")
	}

	// Create member with default role (member)
	now := time.Now()
	member := &entity.StartupMember{
		ID:        uuid.New().String(),
		StartupID: startup.ID,
		UserID:    userID,
		Role:      entity.MemberRoleMember,
		Status:    entity.MemberStatusActive,
		JoinedAt:  &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if existing != nil {
		// Update existing member
		member.ID = existing.ID
		member.CreatedAt = existing.CreatedAt
		if err := uc.memberRepo.Update(ctx, member); err != nil {
			return nil, err
		}
	} else {
		// Create new member
		if err := uc.memberRepo.Create(ctx, member); err != nil {
			return nil, err
		}
	}

	return &dto.MemberOutput{
		ID:        member.ID,
		UserID:    member.UserID,
		StartupID: member.StartupID,
		Role:      string(member.Role),
		Status:    string(member.Status),
		User: dto.UserSummary{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		JoinedAt: member.JoinedAt,
	}, nil
}




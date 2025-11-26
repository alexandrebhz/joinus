package member

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type RequestJoinUseCase struct {
	startupRepo    repository.StartupRepository
	memberRepo     repository.StartupMemberRepository
	userRepo       repository.UserRepository
	notifyService  port.NotificationService
	logger         logger.Logger
}

func NewRequestJoinUseCase(
	startupRepo repository.StartupRepository,
	memberRepo repository.StartupMemberRepository,
	userRepo repository.UserRepository,
	notifyService port.NotificationService,
	logger logger.Logger,
) *RequestJoinUseCase {
	return &RequestJoinUseCase{
		startupRepo:   startupRepo,
		memberRepo:    memberRepo,
		userRepo:      userRepo,
		notifyService: notifyService,
		logger:        logger,
	}
}

func (uc *RequestJoinUseCase) Execute(ctx context.Context, input dto.JoinRequestInput, userID string) (*dto.MemberOutput, error) {
	// Find startup by slug
	startup, err := uc.startupRepo.FindBySlug(ctx, input.StartupSlug)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

	// Get user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}

	// Check if already a member
	existing, _ := uc.memberRepo.FindByUserAndStartup(ctx, userID, startup.ID)
	if existing != nil {
		if existing.Status == entity.MemberStatusActive {
			return nil, errors.NewBadRequestError("you are already a member of this startup")
		}
		if existing.Status == entity.MemberStatusPending {
			return nil, errors.NewBadRequestError("you already have a pending request")
		}
	}

	// Create pending member record
	now := time.Now()
	member := &entity.StartupMember{
		ID:        uuid.New().String(),
		StartupID: startup.ID,
		UserID:    userID,
		Role:      entity.MemberRoleMember,
		Status:    entity.MemberStatusPending,
		JoinedAt:  nil,
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

	// Notify startup owners/admins
	if err := uc.notifyService.NotifyJoinRequest(ctx, member, startup); err != nil {
		uc.logger.Warn("Failed to send join request notification: %v", err)
		// Don't fail the request if notification fails
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




package startup

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/pkg/logger"
	"github.com/startup-job-board/backend/pkg/utils"
)

type CreateStartupUseCase struct {
	startupRepo repository.StartupRepository
	memberRepo  repository.StartupMemberRepository
	tokenGen    port.TokenService
	logger      logger.Logger
}

func NewCreateStartupUseCase(
	startupRepo repository.StartupRepository,
	memberRepo repository.StartupMemberRepository,
	tokenGen port.TokenService,
	logger logger.Logger,
) *CreateStartupUseCase {
	return &CreateStartupUseCase{
		startupRepo: startupRepo,
		memberRepo:  memberRepo,
		tokenGen:    tokenGen,
		logger:      logger,
	}
}

func (uc *CreateStartupUseCase) Execute(ctx context.Context, input dto.CreateStartupInput, userID string) (*dto.StartupOutput, error) {
	// Generate API token
	tokenStr, err := uc.tokenGen.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Generate slug
	baseSlug := utils.GenerateSlug(input.Name)
	slug := utils.GenerateUniqueSlug(baseSlug, func(s string) bool {
		_, err := uc.startupRepo.FindBySlug(ctx, s)
		return err == nil
	})

	// Create startup entity
	startup := &entity.Startup{
		ID:              uuid.New().String(),
		Name:            input.Name,
		Slug:            slug,
		Description:     input.Description,
		Website:         input.Website,
		FoundedYear:     input.FoundedYear,
		Industry:        input.Industry,
		CompanySize:     input.CompanySize,
		Location:        input.Location,
		APIToken:        tokenStr,
		AllowPublicJoin: input.AllowPublicJoin,
		Status:          entity.StartupStatusActive,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.startupRepo.Create(ctx, startup); err != nil {
		return nil, err
	}

	// Create member record for creator as owner
	member := &entity.StartupMember{
		ID:        uuid.New().String(),
		StartupID: startup.ID,
		UserID:    userID,
		Role:      entity.MemberRoleOwner,
		Status:    entity.MemberStatusActive,
		InvitedAt: time.Now(),
		JoinedAt:  &[]time.Time{time.Now()}[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.memberRepo.Create(ctx, member); err != nil {
		// Rollback startup creation
		uc.startupRepo.Delete(ctx, startup.ID)
		return nil, err
	}

	return uc.toOutput(startup), nil
}

func (uc *CreateStartupUseCase) toOutput(startup *entity.Startup) *dto.StartupOutput {
	return &dto.StartupOutput{
		ID:              startup.ID,
		Name:            startup.Name,
		Slug:            startup.Slug,
		Description:     startup.Description,
		LogoURL:         startup.LogoURL,
		Website:         startup.Website,
		FoundedYear:     startup.FoundedYear,
		Industry:        startup.Industry,
		CompanySize:     startup.CompanySize,
		Location:        startup.Location,
		AllowPublicJoin: startup.AllowPublicJoin,
		Status:          string(startup.Status),
		CreatedAt:       startup.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       startup.UpdatedAt.Format(time.RFC3339),
	}
}




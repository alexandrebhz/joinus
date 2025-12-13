package contact

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/logger"
)

type CreateContactUseCase struct {
	contactRepo repository.ContactRepository
	logger      logger.Logger
}

func NewCreateContactUseCase(
	contactRepo repository.ContactRepository,
	logger logger.Logger,
) *CreateContactUseCase {
	return &CreateContactUseCase{
		contactRepo: contactRepo,
		logger:      logger,
	}
}

func (uc *CreateContactUseCase) Execute(ctx context.Context, input dto.CreateContactInput) (*dto.ContactOutput, error) {
	now := time.Now()

	contact := &entity.Contact{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Email:     input.Email,
		Subject:   input.Subject,
		Message:   input.Message,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.contactRepo.Create(ctx, contact); err != nil {
		uc.logger.Error("Failed to create contact: %v", err)
		return nil, err
	}

	return &dto.ContactOutput{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		Subject:   contact.Subject,
		Message:   contact.Message,
		CreatedAt: contact.CreatedAt.Format(time.RFC3339),
	}, nil
}

package postgres

import (
	"context"

	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type ContactRepositoryImpl struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) repository.ContactRepository {
	return &ContactRepositoryImpl{db: db}
}

func (r *ContactRepositoryImpl) Create(ctx context.Context, contact *entity.Contact) error {
	model := r.toModel(contact)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ContactRepositoryImpl) toModel(contact *entity.Contact) *gorm_model.Contact {
	return &gorm_model.Contact{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		Subject:   contact.Subject,
		Message:   contact.Message,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}

func (r *ContactRepositoryImpl) toDomain(model *gorm_model.Contact) *entity.Contact {
	return &entity.Contact{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Subject:   model.Subject,
		Message:   model.Message,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

package postgres

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type InvitationRepositoryImpl struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) repository.InvitationRepository {
	return &InvitationRepositoryImpl{db: db}
}

func (r *InvitationRepositoryImpl) Create(ctx context.Context, invitation *entity.Invitation) error {
	model := r.toModel(invitation)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *InvitationRepositoryImpl) Update(ctx context.Context, invitation *entity.Invitation) error {
	model := r.toModel(invitation)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *InvitationRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Invitation, error) {
	var model gorm_model.Invitation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *InvitationRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.Invitation, error) {
	var model gorm_model.Invitation
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *InvitationRepositoryImpl) FindByEmail(ctx context.Context, email string) ([]*entity.Invitation, error) {
	var models []gorm_model.Invitation
	if err := r.db.WithContext(ctx).Where("email = ?", email).Find(&models).Error; err != nil {
		return nil, err
	}

	invitations := make([]*entity.Invitation, len(models))
	for i, m := range models {
		invitations[i] = r.toDomain(&m)
	}
	return invitations, nil
}

func (r *InvitationRepositoryImpl) FindPendingByStartup(ctx context.Context, startupID string) ([]*entity.Invitation, error) {
	var models []gorm_model.Invitation
	if err := r.db.WithContext(ctx).Where("startup_id = ? AND status = ?", startupID, "pending").Find(&models).Error; err != nil {
		return nil, err
	}

	invitations := make([]*entity.Invitation, len(models))
	for i, m := range models {
		invitations[i] = r.toDomain(&m)
	}
	return invitations, nil
}

func (r *InvitationRepositoryImpl) ExpireOldInvitations(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&gorm_model.Invitation{}).
		Where("status = ? AND expires_at < ?", "pending", now).
		Update("status", "expired").Error
}

func (r *InvitationRepositoryImpl) toModel(invitation *entity.Invitation) *gorm_model.Invitation {
	return &gorm_model.Invitation{
		ID:        invitation.ID,
		StartupID: invitation.StartupID,
		Email:     invitation.Email,
		Token:     invitation.Token,
		Role:      string(invitation.Role),
		InvitedBy: invitation.InvitedBy,
		ExpiresAt: invitation.ExpiresAt,
		AcceptedAt: invitation.AcceptedAt,
		Status:    string(invitation.Status),
		CreatedAt: invitation.CreatedAt,
	}
}

func (r *InvitationRepositoryImpl) toDomain(model *gorm_model.Invitation) *entity.Invitation {
	return &entity.Invitation{
		ID:        model.ID,
		StartupID: model.StartupID,
		Email:     model.Email,
		Token:     model.Token,
		Role:      entity.MemberRole(model.Role),
		InvitedBy: model.InvitedBy,
		ExpiresAt: model.ExpiresAt,
		AcceptedAt: model.AcceptedAt,
		Status:    entity.InvitationStatus(model.Status),
		CreatedAt: model.CreatedAt,
	}
}




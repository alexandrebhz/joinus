package postgres

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type StartupMemberRepositoryImpl struct {
	db *gorm.DB
}

func NewStartupMemberRepository(db *gorm.DB) repository.StartupMemberRepository {
	return &StartupMemberRepositoryImpl{db: db}
}

func (r *StartupMemberRepositoryImpl) Create(ctx context.Context, member *entity.StartupMember) error {
	model := r.toModel(member)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *StartupMemberRepositoryImpl) Update(ctx context.Context, member *entity.StartupMember) error {
	model := r.toModel(member)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *StartupMemberRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&gorm_model.StartupMember{}, id).Error
}

func (r *StartupMemberRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.StartupMember, error) {
	var model gorm_model.StartupMember
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *StartupMemberRepositoryImpl) FindByUserAndStartup(ctx context.Context, userID, startupID string) (*entity.StartupMember, error) {
	var model gorm_model.StartupMember
	if err := r.db.WithContext(ctx).Where("user_id = ? AND startup_id = ?", userID, startupID).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *StartupMemberRepositoryImpl) FindByStartupID(ctx context.Context, startupID string) ([]*entity.StartupMember, error) {
	var models []gorm_model.StartupMember
	if err := r.db.WithContext(ctx).Where("startup_id = ?", startupID).Find(&models).Error; err != nil {
		return nil, err
	}

	members := make([]*entity.StartupMember, len(models))
	for i, m := range models {
		members[i] = r.toDomain(&m)
	}
	return members, nil
}

func (r *StartupMemberRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*entity.StartupMember, error) {
	var models []gorm_model.StartupMember
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	members := make([]*entity.StartupMember, len(models))
	for i, m := range models {
		members[i] = r.toDomain(&m)
	}
	return members, nil
}

func (r *StartupMemberRepositoryImpl) toModel(member *entity.StartupMember) *gorm_model.StartupMember {
	return &gorm_model.StartupMember{
		ID:        member.ID,
		StartupID: member.StartupID,
		UserID:    member.UserID,
		Role:      string(member.Role),
		Status:    string(member.Status),
		InvitedBy: member.InvitedBy,
		InvitedAt: member.InvitedAt,
		JoinedAt:  member.JoinedAt,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func (r *StartupMemberRepositoryImpl) toDomain(model *gorm_model.StartupMember) *entity.StartupMember {
	return &entity.StartupMember{
		ID:        model.ID,
		StartupID: model.StartupID,
		UserID:    model.UserID,
		Role:      entity.MemberRole(model.Role),
		Status:    entity.MemberStatus(model.Status),
		InvitedBy: model.InvitedBy,
		InvitedAt: model.InvitedAt,
		JoinedAt:  model.JoinedAt,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}




package postgres

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	model := r.toModel(user)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	model := r.toModel(user)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var model gorm_model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model gorm_model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, page, pageSize int, search string) ([]*entity.User, int64, error) {
	q := r.db.WithContext(ctx).Model(&gorm_model.User{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("email ILIKE ? OR name ILIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var models []gorm_model.User
	if err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*entity.User, len(models))
	for i := range models {
		out[i] = r.toDomain(&models[i])
	}
	return out, total, nil
}

func (r *UserRepositoryImpl) toModel(user *entity.User) *gorm_model.User {
	return &gorm_model.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		Name:      user.Name,
		Role:      string(user.Role),
		StartupID: user.StartupID,
		Status:       string(user.Status),
		TokenVersion: user.TokenVersion,
		CreatedAt:    user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (r *UserRepositoryImpl) toDomain(model *gorm_model.User) *entity.User {
	return &entity.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		Name:      model.Name,
		Role:      entity.UserRole(model.Role),
		StartupID: model.StartupID,
		Status:       entity.UserStatus(model.Status),
		TokenVersion: model.TokenVersion,
		CreatedAt:    model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}




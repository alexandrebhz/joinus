package postgres

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type OAuthLoginCodeRepositoryImpl struct {
	db *gorm.DB
}

func NewOAuthLoginCodeRepository(db *gorm.DB) repository.OAuthLoginCodeRepository {
	return &OAuthLoginCodeRepositoryImpl{db: db}
}

func (r *OAuthLoginCodeRepositoryImpl) Create(ctx context.Context, code *entity.OAuthLoginCode) error {
	return r.db.WithContext(ctx).Create(&gorm_model.OAuthLoginCode{
		ID: code.ID, CodeHash: code.CodeHash,
		AccessToken: code.AccessToken, RefreshToken: code.RefreshToken,
		UserID: code.UserID, ExpiresAt: code.ExpiresAt, UsedAt: code.UsedAt,
		CreatedAt: code.CreatedAt,
	}).Error
}

func (r *OAuthLoginCodeRepositoryImpl) FindByCodeHash(ctx context.Context, codeHash string) (*entity.OAuthLoginCode, error) {
	var m gorm_model.OAuthLoginCode
	if err := r.db.WithContext(ctx).Where("code_hash = ?", codeHash).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity.OAuthLoginCode{
		ID: m.ID, CodeHash: m.CodeHash,
		AccessToken: m.AccessToken, RefreshToken: m.RefreshToken,
		UserID: m.UserID, ExpiresAt: m.ExpiresAt, UsedAt: m.UsedAt,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (r *OAuthLoginCodeRepositoryImpl) MarkUsed(ctx context.Context, id string, usedAt time.Time) error {
	res := r.db.WithContext(ctx).Model(&gorm_model.OAuthLoginCode{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", usedAt)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OAuthLoginCodeRepositoryImpl) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ? OR used_at IS NOT NULL", time.Now().Add(-1*time.Hour)).
		Delete(&gorm_model.OAuthLoginCode{}).Error
}

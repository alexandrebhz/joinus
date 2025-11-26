package postgres

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type FileRepositoryImpl struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) repository.FileRepository {
	return &FileRepositoryImpl{db: db}
}

func (r *FileRepositoryImpl) Create(ctx context.Context, file *entity.File) error {
	model := r.toModel(file)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *FileRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.File, error) {
	var model gorm_model.File
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *FileRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&gorm_model.File{}, id).Error
}

func (r *FileRepositoryImpl) toModel(file *entity.File) *gorm_model.File {
	return &gorm_model.File{
		ID:         file.ID,
		FileName:   file.FileName,
		FileSize:   file.FileSize,
		MimeType:   file.MimeType,
		StorageKey: file.StorageKey,
		URL:        file.URL,
		UploadedBy: file.UploadedBy,
		CreatedAt:  file.CreatedAt,
	}
}

func (r *FileRepositoryImpl) toDomain(model *gorm_model.File) *entity.File {
	return &entity.File{
		ID:         model.ID,
		FileName:   model.FileName,
		FileSize:   model.FileSize,
		MimeType:   model.MimeType,
		StorageKey: model.StorageKey,
		URL:        model.URL,
		UploadedBy: model.UploadedBy,
		CreatedAt:  model.CreatedAt,
	}
}




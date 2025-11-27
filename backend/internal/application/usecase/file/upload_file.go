package file

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
	"path/filepath"
	"strings"
)

const (
	MaxFileSize = 2 * 1024 * 1024 // 2MB
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type UploadFileUseCase struct {
	fileRepo       repository.FileRepository
	storageService port.StorageService
	logger         logger.Logger
}

func NewUploadFileUseCase(
	fileRepo repository.FileRepository,
	storageService port.StorageService,
	logger logger.Logger,
) *UploadFileUseCase {
	return &UploadFileUseCase{
		fileRepo:       fileRepo,
		storageService: storageService,
		logger:         logger,
	}
}

func (uc *UploadFileUseCase) Execute(ctx context.Context, fileData []byte, fileName string, mimeType string, userID string) (*dto.FileOutput, error) {
	// Check if storage service is available
	if uc.storageService == nil {
		return nil, errors.NewBadRequestError("file storage is not configured")
	}

	// Validate file size
	if len(fileData) > MaxFileSize {
		return nil, errors.NewBadRequestError("file size exceeds 2MB limit")
	}

	// Validate MIME type
	if !allowedMimeTypes[strings.ToLower(mimeType)] {
		return nil, errors.NewBadRequestError("only image files are allowed (JPEG, PNG, GIF, WebP)")
	}

	// Generate storage key
	ext := filepath.Ext(fileName)
	storageKey := "uploads/" + uuid.New().String() + ext

	// Upload to storage
	url, err := uc.storageService.Upload(ctx, fileData, storageKey, mimeType)
	if err != nil {
		return nil, err
	}

	// Create file record
	file := &entity.File{
		ID:         uuid.New().String(),
		FileName:   fileName,
		FileSize:   int64(len(fileData)),
		MimeType:   mimeType,
		StorageKey: storageKey,
		URL:        url,
		UploadedBy: userID,
		CreatedAt:  time.Now(),
	}

	if err := uc.fileRepo.Create(ctx, file); err != nil {
		// Try to delete from storage if DB insert fails
		if uc.storageService != nil {
			uc.storageService.Delete(ctx, storageKey)
		}
		return nil, err
	}

	return &dto.FileOutput{
		ID:        file.ID,
		FileName:  file.FileName,
		FileSize:  file.FileSize,
		MimeType:  file.MimeType,
		URL:       file.URL,
		CreatedAt: file.CreatedAt.Format(time.RFC3339),
	}, nil
}


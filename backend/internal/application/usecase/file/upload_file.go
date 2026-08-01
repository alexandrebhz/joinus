package file

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

const (
	MaxFileSize = 2 * 1024 * 1024 // 2MB
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var extensionMimeFallback = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
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
	if uc.storageService == nil {
		return nil, errors.NewBadRequestError("file storage is not configured")
	}

	if len(fileData) > MaxFileSize {
		return nil, errors.NewBadRequestError("file size exceeds 2MB limit")
	}

	resolvedMime, err := resolveImageMimeType(fileData, fileName, mimeType)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(fileName)
	storageKey := "uploads/" + uuid.New().String() + ext

	url, err := uc.storageService.Upload(ctx, fileData, storageKey, resolvedMime)
	if err != nil {
		return nil, err
	}

	file := &entity.File{
		ID:         uuid.New().String(),
		FileName:   fileName,
		FileSize:   int64(len(fileData)),
		MimeType:   resolvedMime,
		StorageKey: storageKey,
		URL:        url,
		UploadedBy: userID,
		CreatedAt:  time.Now(),
	}

	if err := uc.fileRepo.Create(ctx, file); err != nil {
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

func resolveImageMimeType(fileData []byte, fileName, clientMime string) (string, error) {
	clientMime = strings.ToLower(strings.TrimSpace(clientMime))
	if strings.Contains(clientMime, "svg") {
		return "", errors.NewBadRequestError("only image files are allowed (JPEG, PNG, GIF, WebP)")
	}

	detected := http.DetectContentType(fileData)
	if strings.HasPrefix(detected, "image/svg") {
		return "", errors.NewBadRequestError("only image files are allowed (JPEG, PNG, GIF, WebP)")
	}

	if allowedMimeTypes[detected] {
		return detected, nil
	}

	if detected == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(fileName))
		if fallback, ok := extensionMimeFallback[ext]; ok {
			return fallback, nil
		}
	}

	return "", errors.NewBadRequestError("only image files are allowed (JPEG, PNG, GIF, WebP)")
}

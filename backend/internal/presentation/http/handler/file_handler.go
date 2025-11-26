package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	fileusecase "github.com/startup-job-board/backend/internal/application/usecase/file"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
)

type FileHandler struct {
	uploadUseCase *fileusecase.UploadFileUseCase
}

func NewFileHandler(uploadUseCase *fileusecase.UploadFileUseCase) *FileHandler {
	return &FileHandler{
		uploadUseCase: uploadUseCase,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	// Read file
	src, err := file.Open()
	if err != nil {
		response.BadRequest(c, "failed to open file")
		return
	}
	defer src.Close()

	// Read file into byte slice
	fileData := make([]byte, file.Size)
	if _, err := src.Read(fileData); err != nil {
		response.BadRequest(c, "failed to read file")
		return
	}

	userID := middleware.GetUserID(c)
	result, err := h.uploadUseCase.Execute(c.Request.Context(), fileData, file.Filename, file.Header.Get("Content-Type"), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	response.Success(c, result)
}




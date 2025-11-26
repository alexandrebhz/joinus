package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/pkg/errors"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

func Error(c *gin.Context, statusCode int, err error) {
	appErr, ok := err.(*errors.AppError)
	if ok {
		c.JSON(statusCode, ErrorResponse{
			Success: false,
			Error:   appErr.Message,
			Code:    appErr.Code,
		})
		return
	}

	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, errors.NewBadRequestError(message))
}

func Unauthorized(c *gin.Context) {
	Error(c, http.StatusUnauthorized, errors.ErrUnauthorized)
}

func Forbidden(c *gin.Context) {
	Error(c, http.StatusForbidden, errors.ErrForbidden)
}

func NotFound(c *gin.Context, resource string) {
	Error(c, http.StatusNotFound, errors.NewNotFoundError(resource))
}

func InternalError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, errors.ErrInternalError)
}




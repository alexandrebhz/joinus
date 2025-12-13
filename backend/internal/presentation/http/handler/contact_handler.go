package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	contactusecase "github.com/startup-job-board/backend/internal/application/usecase/contact"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
)

type ContactHandler struct {
	createContactUseCase *contactusecase.CreateContactUseCase
	validator            *validator.Validator
}

func NewContactHandler(
	createContactUseCase *contactusecase.CreateContactUseCase,
	validator *validator.Validator,
) *ContactHandler {
	return &ContactHandler{
		createContactUseCase: createContactUseCase,
		validator:            validator,
	}
}

func (h *ContactHandler) Create(c *gin.Context) {
	var input dto.CreateContactInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	contact, err := h.createContactUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	response.Success(c, contact)
}

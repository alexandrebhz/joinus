package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	authusecase "github.com/startup-job-board/backend/internal/application/usecase/auth"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
)

type AuthHandler struct {
	registerUseCase    *authusecase.RegisterUseCase
	loginUseCase       *authusecase.LoginUseCase
	refreshTokenUseCase *authusecase.RefreshTokenUseCase
	validator          *validator.Validator
}

func NewAuthHandler(
	registerUseCase *authusecase.RegisterUseCase,
	loginUseCase *authusecase.LoginUseCase,
	refreshTokenUseCase *authusecase.RefreshTokenUseCase,
	validator *validator.Validator,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase:    registerUseCase,
		loginUseCase:       loginUseCase,
		refreshTokenUseCase: refreshTokenUseCase,
		validator:          validator,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.registerUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	response.Success(c, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.loginUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err)
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input dto.RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.refreshTokenUseCase.Execute(c.Request.Context(), input.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err)
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	// This would need a GetUser use case
	response.Success(c, gin.H{"user_id": userID})
}




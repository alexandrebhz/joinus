package handler

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	authusecase "github.com/startup-job-board/backend/internal/application/usecase/auth"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
	"github.com/startup-job-board/backend/pkg/errors"
)

type AuthHandler struct {
	registerUseCase     *authusecase.RegisterUseCase
	loginUseCase        *authusecase.LoginUseCase
	refreshTokenUseCase *authusecase.RefreshTokenUseCase
	getMeUseCase        *authusecase.GetMeUseCase
	startOAuthUseCase   *authusecase.StartOAuthUseCase
	completeOAuthUseCase *authusecase.CompleteOAuthUseCase
	frontendURL         string
	validator           *validator.Validator
}

func NewAuthHandler(
	registerUseCase *authusecase.RegisterUseCase,
	loginUseCase *authusecase.LoginUseCase,
	refreshTokenUseCase *authusecase.RefreshTokenUseCase,
	getMeUseCase *authusecase.GetMeUseCase,
	startOAuthUseCase *authusecase.StartOAuthUseCase,
	completeOAuthUseCase *authusecase.CompleteOAuthUseCase,
	frontendURL string,
	validator *validator.Validator,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase, loginUseCase: loginUseCase,
		refreshTokenUseCase: refreshTokenUseCase, getMeUseCase: getMeUseCase,
		startOAuthUseCase: startOAuthUseCase, completeOAuthUseCase: completeOAuthUseCase,
		frontendURL: frontendURL, validator: validator,
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
	result, err := h.registerUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, result)
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
	result, err := h.getMeUseCase.Execute(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == "NOT_FOUND" {
			response.Error(c, http.StatusNotFound, err)
			return
		}
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) OAuthStart(c *gin.Context) {
	provider := c.Param("provider")
	authURL, state, err := h.startOAuthUseCase.Execute(provider)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == "BAD_REQUEST" {
			// Apple/GitHub stubs
			response.Error(c, http.StatusNotImplemented, err)
			return
		}
		mapUCError(c, err)
		return
	}
	if authURL == "" {
		response.Error(c, http.StatusNotImplemented, errors.NewBadRequestError(provider+" oauth is not implemented yet"))
		return
	}
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)
	c.Redirect(http.StatusFound, authURL)
}

func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	cookieState, _ := c.Cookie("oauth_state")
	if code == "" || state == "" || cookieState == "" || state != cookieState {
		response.Error(c, http.StatusBadRequest, errors.NewBadRequestError("invalid oauth state"))
		return
	}
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	result, err := h.completeOAuthUseCase.Execute(c.Request.Context(), provider, code)
	if err != nil {
		mapUCError(c, err)
		return
	}

	// Redirect to frontend with tokens in fragment would be ideal; query is used for SPA pickup.
	redirect := h.frontendURL + "/auth/callback?access_token=" + url.QueryEscape(result.AccessToken) +
		"&refresh_token=" + url.QueryEscape(result.RefreshToken)
	c.Redirect(http.StatusFound, redirect)
}

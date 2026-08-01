package handler

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	authusecase "github.com/startup-job-board/backend/internal/application/usecase/auth"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
	"github.com/startup-job-board/backend/pkg/errors"
)

type AuthHandler struct {
	registerUseCase          *authusecase.RegisterUseCase
	loginUseCase             *authusecase.LoginUseCase
	refreshTokenUseCase      *authusecase.RefreshTokenUseCase
	logoutUseCase            *authusecase.LogoutUseCase
	getMeUseCase             *authusecase.GetMeUseCase
	startOAuthUseCase        *authusecase.StartOAuthUseCase
	completeOAuthUseCase     *authusecase.CompleteOAuthUseCase
	issueLoginCodeUseCase    *authusecase.IssueOAuthLoginCodeUseCase
	exchangeLoginCodeUseCase *authusecase.ExchangeOAuthLoginCodeUseCase
	frontendURL              string
	secureCookies            bool
	validator                *validator.Validator
}

func NewAuthHandler(
	registerUseCase *authusecase.RegisterUseCase,
	loginUseCase *authusecase.LoginUseCase,
	refreshTokenUseCase *authusecase.RefreshTokenUseCase,
	logoutUseCase *authusecase.LogoutUseCase,
	getMeUseCase *authusecase.GetMeUseCase,
	startOAuthUseCase *authusecase.StartOAuthUseCase,
	completeOAuthUseCase *authusecase.CompleteOAuthUseCase,
	issueLoginCodeUseCase *authusecase.IssueOAuthLoginCodeUseCase,
	exchangeLoginCodeUseCase *authusecase.ExchangeOAuthLoginCodeUseCase,
	frontendURL string,
	secureCookies bool,
	validator *validator.Validator,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase, loginUseCase: loginUseCase,
		refreshTokenUseCase: refreshTokenUseCase, logoutUseCase: logoutUseCase, getMeUseCase: getMeUseCase,
		startOAuthUseCase: startOAuthUseCase, completeOAuthUseCase: completeOAuthUseCase,
		issueLoginCodeUseCase: issueLoginCodeUseCase, exchangeLoginCodeUseCase: exchangeLoginCodeUseCase,
		frontendURL: frontendURL, secureCookies: secureCookies, validator: validator,
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

func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.logoutUseCase.Execute(c.Request.Context(), middleware.GetUserID(c)); err != nil {
		response.Error(c, http.StatusUnauthorized, err)
		return
	}
	response.Success(c, gin.H{"message": "logged out successfully"})
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
	h.setOAuthStateCookie(c, state, 600)
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
	h.clearOAuthStateCookie(c)

	result, err := h.completeOAuthUseCase.Execute(c.Request.Context(), provider, code)
	if err != nil {
		mapUCError(c, err)
		return
	}

	loginCode, err := h.issueLoginCodeUseCase.Execute(c.Request.Context(), result)
	if err != nil {
		mapUCError(c, err)
		return
	}

	// Only an opaque one-time code is sent to the frontend — never JWTs in the URL.
	redirect := strings.TrimRight(h.frontendURL, "/") + "/auth/callback?code=" + url.QueryEscape(loginCode)
	c.Redirect(http.StatusFound, redirect)
}

func (h *AuthHandler) ExchangeOAuthCode(c *gin.Context) {
	var input dto.ExchangeOAuthCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.exchangeLoginCodeUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) setOAuthStateCookie(c *gin.Context, state string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, maxAge, "/", "", h.secureCookies, true)
}

func (h *AuthHandler) clearOAuthStateCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", "", -1, "/", "", h.secureCookies, true)
}

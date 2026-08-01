package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	startupusecase "github.com/startup-job-board/backend/internal/application/usecase/startup"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/utils"
)

type StartupHandler struct {
	createUseCase *startupusecase.CreateStartupUseCase
	updateUseCase *startupusecase.UpdateStartupUseCase
	getUseCase    *startupusecase.GetStartupUseCase
	listUseCase   *startupusecase.ListStartupsUseCase
	validator     *validator.Validator
}

func NewStartupHandler(
	createUseCase *startupusecase.CreateStartupUseCase,
	updateUseCase *startupusecase.UpdateStartupUseCase,
	getUseCase *startupusecase.GetStartupUseCase,
	listUseCase *startupusecase.ListStartupsUseCase,
	validator *validator.Validator,
) *StartupHandler {
	return &StartupHandler{
		createUseCase: createUseCase,
		updateUseCase: updateUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		validator:     validator,
	}
}

func (h *StartupHandler) Create(c *gin.Context) {
	var input dto.CreateStartupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	userRoleStr := middleware.GetUserRole(c)
	userRole := entity.UserRole(userRoleStr)

	result, err := h.createUseCase.Execute(c.Request.Context(), input, userID, userRole)
	if err != nil {
		// Check error type to return appropriate status code
		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case "FORBIDDEN":
				response.Error(c, http.StatusForbidden, err)
			case "NOT_FOUND":
				response.Error(c, http.StatusNotFound, err)
			default:
				response.Error(c, http.StatusBadRequest, err)
			}
		} else {
			response.Error(c, http.StatusBadRequest, err)
		}
		return
	}

	response.Success(c, result)
}

func (h *StartupHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input dto.UpdateStartupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input.ID = id
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	result, err := h.updateUseCase.Execute(c.Request.Context(), input, userID)
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

func (h *StartupHandler) Get(c *gin.Context) {
	id := c.Param("id")
	result, err := h.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err)
		return
	}

	response.Success(c, result)
}

func (h *StartupHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	result, err := h.getUseCase.ExecuteBySlug(c.Request.Context(), slug)
	if err != nil {
		response.Error(c, http.StatusNotFound, err)
		return
	}

	response.Success(c, result)
}

func (h *StartupHandler) List(c *gin.Context) {
	filter := repository.StartupFilter{}
	authenticated := middleware.GetUserID(c) != ""
	trusted := middleware.IsInternalTrusted(c)

	if industry := c.Query("industry"); industry != "" {
		filter.Industry = industry
	}
	if status := c.Query("status"); status != "" {
		filter.Status = entity.StartupStatus(status)
	}
	// Unauthenticated scrapers cannot enumerate non-active startups.
	if !authenticated && !trusted {
		filter.Status = entity.StartupStatusActive
	}
	if search := c.Query("search"); search != "" {
		filter.Search = search
	}
	if location := c.Query("location"); location != "" {
		filter.Location = location
	}
	if companySize := c.Query("company_size"); companySize != "" {
		filter.CompanySize = companySize
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	maxPageSize := utils.MaxPageSizePublic
	if authenticated {
		maxPageSize = utils.MaxPageSizeAuth
	}
	page, pageSize = utils.ClampPagination(page, pageSize, maxPageSize)
	filter.Page = page
	filter.PageSize = pageSize
	filter.OrderBy = c.DefaultQuery("order_by", "created_at")
	filter.OrderDir = c.DefaultQuery("order_dir", "DESC")

	startups, total, err := h.listUseCase.Execute(c.Request.Context(), filter, true)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: utils.CalculateTotalPages(total, pageSize),
	}

	response.SuccessWithMeta(c, startups, meta)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	startupusecase "github.com/startup-job-board/backend/internal/application/usecase/startup"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
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
	result, err := h.createUseCase.Execute(c.Request.Context(), input, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
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

	result, err := h.updateUseCase.Execute(c.Request.Context(), input)
	if err != nil {
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

	if industry := c.Query("industry"); industry != "" {
		filter.Industry = industry
	}
	if status := c.Query("status"); status != "" {
		filter.Status = entity.StartupStatus(status)
	}
	if search := c.Query("search"); search != "" {
		filter.Search = search
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	filter.Page = page
	filter.PageSize = pageSize
	filter.OrderBy = c.DefaultQuery("order_by", "created_at")
	filter.OrderDir = c.DefaultQuery("order_dir", "DESC")

	startups, total, err := h.listUseCase.Execute(c.Request.Context(), filter)
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


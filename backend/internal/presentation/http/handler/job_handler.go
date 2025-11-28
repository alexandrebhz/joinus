package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	jobusecase "github.com/startup-job-board/backend/internal/application/usecase/job"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
	"github.com/startup-job-board/backend/pkg/utils"
)

type JobHandler struct {
	createUseCase *jobusecase.CreateJobUseCase
	updateUseCase *jobusecase.UpdateJobUseCase
	listUseCase   *jobusecase.ListJobsUseCase
	deleteUseCase *jobusecase.DeleteJobUseCase
	jobRepo       repository.JobRepository
	startupRepo   repository.StartupRepository
	validator     *validator.Validator
}

func NewJobHandler(
	createUseCase *jobusecase.CreateJobUseCase,
	updateUseCase *jobusecase.UpdateJobUseCase,
	listUseCase *jobusecase.ListJobsUseCase,
	deleteUseCase *jobusecase.DeleteJobUseCase,
	jobRepo repository.JobRepository,
	startupRepo repository.StartupRepository,
	validator *validator.Validator,
) *JobHandler {
	return &JobHandler{
		createUseCase: createUseCase,
		updateUseCase: updateUseCase,
		listUseCase:   listUseCase,
		deleteUseCase: deleteUseCase,
		jobRepo:       jobRepo,
		startupRepo:   startupRepo,
		validator:     validator,
	}
}

func (h *JobHandler) Create(c *gin.Context) {
	var input dto.CreateJobInput
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

func (h *JobHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input dto.UpdateJobInput
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
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	response.Success(c, result)
}

func (h *JobHandler) List(c *gin.Context) {
	filter := repository.JobFilter{}

	if startupID := c.Query("startup_id"); startupID != "" {
		filter.StartupID = startupID
	}
	if jobType := c.Query("job_type"); jobType != "" {
		filter.JobType = entity.JobType(jobType)
	}
	if locationType := c.Query("location_type"); locationType != "" {
		filter.LocationType = entity.LocationType(locationType)
	}
	if status := c.Query("status"); status != "" {
		filter.Status = entity.JobStatus(status)
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

	jobs, total, err := h.listUseCase.Execute(c.Request.Context(), filter)
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

	response.SuccessWithMeta(c, jobs, meta)
}

func (h *JobHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "job id is required")
		return
	}

	job, err := h.jobRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err)
		return
	}

	// Get startup name
	startup, _ := h.startupRepo.FindByID(c.Request.Context(), job.StartupID)
	startupName := ""
	if startup != nil {
		startupName = startup.Name
	}

	// Convert to output DTO
	output := &dto.JobOutput{
		ID:               job.ID,
		StartupID:        job.StartupID,
		StartupName:      startupName,
		Title:            job.Title,
		Description:      job.Description,
		Requirements:     job.Requirements,
		JobType:          string(job.JobType),
		LocationType:     string(job.LocationType),
		City:             job.City,
		Country:          job.Country,
		SalaryMin:        job.SalaryMin,
		SalaryMax:        job.SalaryMax,
		Currency:         job.Currency,
		ApplicationURL:   job.ApplicationURL,
		ApplicationEmail: job.ApplicationEmail,
		Status:           string(job.Status),
		CreatedAt:        job.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        job.UpdatedAt.Format(time.RFC3339),
	}

	if job.ExpiresAt != nil {
		expiresAtStr := job.ExpiresAt.Format(time.RFC3339)
		output.ExpiresAt = &expiresAtStr
	}

	response.Success(c, output)
}

func (h *JobHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.deleteUseCase.Execute(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	response.Success(c, gin.H{"message": "job deleted successfully"})
}

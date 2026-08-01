package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	billingusecase "github.com/startup-job-board/backend/internal/application/usecase/billing"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
)

type BillingHandler struct {
	createCheckoutUseCase *billingusecase.CreateCheckoutUseCase
	handleWebhookUseCase  *billingusecase.HandleWebhookUseCase
	startupRepo           repository.StartupRepository
	authService           *service.AuthorizationService
	validator             *validator.Validator
}

func NewBillingHandler(
	createCheckoutUseCase *billingusecase.CreateCheckoutUseCase,
	handleWebhookUseCase *billingusecase.HandleWebhookUseCase,
	startupRepo repository.StartupRepository,
	authService *service.AuthorizationService,
	validator *validator.Validator,
) *BillingHandler {
	return &BillingHandler{
		createCheckoutUseCase: createCheckoutUseCase,
		handleWebhookUseCase:  handleWebhookUseCase,
		startupRepo:           startupRepo,
		authService:           authService,
		validator:             validator,
	}
}

func (h *BillingHandler) CreateCheckout(c *gin.Context) {
	var input dto.CreateCheckoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	result, err := h.createCheckoutUseCase.Execute(c.Request.Context(), input, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	response.Success(c, result)
}

// Webhook handles POST /billing/webhook. It must receive the raw request body
// so the Stripe signature (Stripe-Signature header) can be verified.
func (h *BillingHandler) Webhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.BadRequest(c, "failed to read request body")
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		response.BadRequest(c, "missing Stripe-Signature header")
		return
	}

	if err := h.handleWebhookUseCase.Execute(c.Request.Context(), payload, signature); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// Status returns the billing plan for every startup the current user manages.
func (h *BillingHandler) Status(c *gin.Context) {
	userID := middleware.GetUserID(c)

	memberships, err := h.authService.GetUserStartups(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	statuses := make([]dto.StartupBillingStatus, 0, len(memberships))
	for _, member := range memberships {
		startup, err := h.startupRepo.FindByID(c.Request.Context(), member.StartupID)
		if err != nil {
			continue
		}

		plan := startup.Plan
		if plan == "" {
			plan = string(entity.StartupPlanFree)
		}

		status := dto.StartupBillingStatus{
			StartupID:   startup.ID,
			StartupName: startup.Name,
			Plan:        plan,
		}
		if startup.PlanExpiresAt != nil {
			expiresAtStr := startup.PlanExpiresAt.Format(time.RFC3339)
			status.PlanExpiresAt = &expiresAtStr
		}
		statuses = append(statuses, status)
	}

	response.Success(c, dto.BillingStatusOutput{Startups: statuses})
}

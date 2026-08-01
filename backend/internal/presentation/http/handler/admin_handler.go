package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	adminusecase "github.com/startup-job-board/backend/internal/application/usecase/admin"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
)

type AdminHandler struct {
	listUsersUC    *adminusecase.ListUsersUseCase
	updateUserUC   *adminusecase.UpdateUserUseCase
	listTeamsUC    *adminusecase.ListTeamsUseCase
	createStartupUC *adminusecase.CreateOrphanStartupUseCase
	linkTeamUC     *adminusecase.LinkStartupTeamUseCase
	validator      *validator.Validator
}

func NewAdminHandler(
	listUsersUC *adminusecase.ListUsersUseCase,
	updateUserUC *adminusecase.UpdateUserUseCase,
	listTeamsUC *adminusecase.ListTeamsUseCase,
	createStartupUC *adminusecase.CreateOrphanStartupUseCase,
	linkTeamUC *adminusecase.LinkStartupTeamUseCase,
	validator *validator.Validator,
) *AdminHandler {
	return &AdminHandler{
		listUsersUC: listUsersUC, updateUserUC: updateUserUC, listTeamsUC: listTeamsUC,
		createStartupUC: createStartupUC, linkTeamUC: linkTeamUC, validator: validator,
	}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	users, total, err := h.listUsersUC.Execute(c.Request.Context(), middleware.GetUserID(c), page, pageSize, c.Query("search"))
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"users": users, "total": total, "page": page, "page_size": pageSize})
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	var input dto.AdminUpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.updateUserUC.Execute(c.Request.Context(), middleware.GetUserID(c), c.Param("id"), input)
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *AdminHandler) ListTeams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	teams, total, err := h.listTeamsUC.Execute(c.Request.Context(), middleware.GetUserID(c), page, pageSize)
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"teams": teams, "total": total})
}

func (h *AdminHandler) CreateStartup(c *gin.Context) {
	var input dto.CreateStartupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.createStartupUC.Execute(c.Request.Context(), middleware.GetUserID(c), input)
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *AdminHandler) LinkStartupTeam(c *gin.Context) {
	var input dto.AdminLinkStartupTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.linkTeamUC.Execute(c.Request.Context(), middleware.GetUserID(c), c.Param("id"), input); err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"status": "ok"})
}

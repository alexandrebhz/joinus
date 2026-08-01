package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/internal/application/dto"
	teamusecase "github.com/startup-job-board/backend/internal/application/usecase/team"
	"github.com/startup-job-board/backend/internal/presentation/http/middleware"
	"github.com/startup-job-board/backend/internal/presentation/http/response"
	"github.com/startup-job-board/backend/internal/presentation/http/validator"
	"github.com/startup-job-board/backend/pkg/errors"
)

type TeamHandler struct {
	createUC          *teamusecase.CreateTeamUseCase
	listUC            *teamusecase.ListMyTeamsUseCase
	getUC             *teamusecase.GetTeamUseCase
	updateUC          *teamusecase.UpdateTeamUseCase
	listMembersUC     *teamusecase.ListMembersUseCase
	inviteUC          *teamusecase.InviteMemberUseCase
	acceptUC          *teamusecase.AcceptInvitationUseCase
	updateMemberUC    *teamusecase.UpdateMemberUseCase
	removeMemberUC    *teamusecase.RemoveMemberUseCase
	listRolesUC       *teamusecase.ListRolesUseCase
	linkStartupUC     *teamusecase.LinkStartupUseCase
	unlinkStartupUC   *teamusecase.UnlinkStartupUseCase
	listStartupsUC    *teamusecase.ListTeamStartupsUseCase
	validator         *validator.Validator
}

func NewTeamHandler(
	createUC *teamusecase.CreateTeamUseCase,
	listUC *teamusecase.ListMyTeamsUseCase,
	getUC *teamusecase.GetTeamUseCase,
	updateUC *teamusecase.UpdateTeamUseCase,
	listMembersUC *teamusecase.ListMembersUseCase,
	inviteUC *teamusecase.InviteMemberUseCase,
	acceptUC *teamusecase.AcceptInvitationUseCase,
	updateMemberUC *teamusecase.UpdateMemberUseCase,
	removeMemberUC *teamusecase.RemoveMemberUseCase,
	listRolesUC *teamusecase.ListRolesUseCase,
	linkStartupUC *teamusecase.LinkStartupUseCase,
	unlinkStartupUC *teamusecase.UnlinkStartupUseCase,
	listStartupsUC *teamusecase.ListTeamStartupsUseCase,
	validator *validator.Validator,
) *TeamHandler {
	return &TeamHandler{
		createUC: createUC, listUC: listUC, getUC: getUC, updateUC: updateUC,
		listMembersUC: listMembersUC, inviteUC: inviteUC, acceptUC: acceptUC,
		updateMemberUC: updateMemberUC, removeMemberUC: removeMemberUC,
		listRolesUC: listRolesUC, linkStartupUC: linkStartupUC,
		unlinkStartupUC: unlinkStartupUC, listStartupsUC: listStartupsUC,
		validator: validator,
	}
}

func mapUCError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		switch appErr.Code {
		case "NOT_FOUND":
			response.Error(c, http.StatusNotFound, err)
		case "FORBIDDEN":
			response.Error(c, http.StatusForbidden, err)
		case "UNAUTHORIZED":
			response.Error(c, http.StatusUnauthorized, err)
		default:
			response.Error(c, http.StatusBadRequest, err)
		}
		return
	}
	response.Error(c, http.StatusBadRequest, err)
}

func (h *TeamHandler) Create(c *gin.Context) {
	var input dto.CreateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.createUC.Execute(c.Request.Context(), input, middleware.GetUserID(c))
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TeamHandler) ListMine(c *gin.Context) {
	result, err := h.listUC.Execute(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TeamHandler) Get(c *gin.Context) {
	result, err := h.getUC.Execute(c.Request.Context(), c.Param("id"), middleware.GetUserID(c))
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TeamHandler) Update(c *gin.Context) {
	var input dto.UpdateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.updateUC.Execute(c.Request.Context(), c.Param("id"), middleware.GetUserID(c), input)
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TeamHandler) ListMembers(c *gin.Context) {
	result, err := h.listMembersUC.Execute(c.Request.Context(), c.Param("id"), middleware.GetUserID(c))
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TeamHandler) Invite(c *gin.Context) {
	var input dto.InviteTeamMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.inviteUC.Execute(c.Request.Context(), c.Param("id"), middleware.GetUserID(c), input); err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"status": "invited"})
}

func (h *TeamHandler) AcceptInvitation(c *gin.Context) {
	var input dto.AcceptTeamInvitationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.acceptUC.Execute(c.Request.Context(), middleware.GetUserID(c), input); err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"status": "accepted"})
}

func (h *TeamHandler) UpdateMember(c *gin.Context) {
	var input dto.UpdateTeamMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.validator.Validate(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.updateMemberUC.Execute(c.Request.Context(), c.Param("id"), c.Param("userId"), middleware.GetUserID(c), input); err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"status": "updated"})
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	if err := h.removeMemberUC.Execute(c.Request.Context(), c.Param("id"), c.Param("userId"), middleware.GetUserID(c)); err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"status": "removed"})
}

func (h *TeamHandler) ListRoles(c *gin.Context) {
	result, err := h.listRolesUC.Execute(c.Request.Context(), c.Param("id"), middleware.GetUserID(c))
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TeamHandler) ListStartups(c *gin.Context) {
	result, err := h.listStartupsUC.Execute(c.Request.Context(), c.Param("id"), middleware.GetUserID(c))
	if err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TeamHandler) LinkStartup(c *gin.Context) {
	if err := h.linkStartupUC.Execute(c.Request.Context(), c.Param("id"), c.Param("startupId"), middleware.GetUserID(c)); err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"status": "linked"})
}

func (h *TeamHandler) UnlinkStartup(c *gin.Context) {
	if err := h.unlinkStartupUC.Execute(c.Request.Context(), c.Param("id"), c.Param("startupId"), middleware.GetUserID(c)); err != nil {
		mapUCError(c, err)
		return
	}
	response.Success(c, gin.H{"status": "unlinked"})
}

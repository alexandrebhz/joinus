package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/startup-job-board/crawler/internal/application/dto"
	"github.com/startup-job-board/crawler/internal/application/usecase/site"
)

// SiteHandler handles HTTP requests for crawl sites
type SiteHandler struct {
	createUseCase *site.CreateSiteUseCase
	updateUseCase *site.UpdateSiteUseCase
	deleteUseCase *site.DeleteSiteUseCase
	getUseCase    *site.GetSiteUseCase
	listUseCase   *site.ListSitesUseCase
	validator     *validator.Validate
}

// NewSiteHandler creates a new SiteHandler
func NewSiteHandler(
	createUseCase *site.CreateSiteUseCase,
	updateUseCase *site.UpdateSiteUseCase,
	deleteUseCase *site.DeleteSiteUseCase,
	getUseCase *site.GetSiteUseCase,
	listUseCase *site.ListSitesUseCase,
) *SiteHandler {
	return &SiteHandler{
		createUseCase: createUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		validator:     validator.New(),
	}
}

// Create handles POST /sites
func (h *SiteHandler) Create(c *gin.Context) {
	var input dto.CreateSiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	site, err := h.createUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": site})
}

// Update handles PUT /sites/:id
func (h *SiteHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input dto.UpdateSiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	site, err := h.updateUseCase.Execute(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": site})
}

// Delete handles DELETE /sites/:id
func (h *SiteHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.deleteUseCase.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site deleted successfully"})
}

// Get handles GET /sites/:id
func (h *SiteHandler) Get(c *gin.Context) {
	id := c.Param("id")
	site, err := h.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": site})
}

// List handles GET /sites
func (h *SiteHandler) List(c *gin.Context) {
	sites, err := h.listUseCase.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sites})
}


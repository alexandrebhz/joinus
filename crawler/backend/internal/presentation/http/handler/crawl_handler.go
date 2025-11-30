package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/crawler/internal/application/usecase/crawl"
)

// CrawlHandler handles HTTP requests for crawl operations
type CrawlHandler struct {
	executeCrawlUseCase *crawl.ExecuteCrawlUseCase
}

// NewCrawlHandler creates a new CrawlHandler
func NewCrawlHandler(executeCrawlUseCase *crawl.ExecuteCrawlUseCase) *CrawlHandler {
	return &CrawlHandler{
		executeCrawlUseCase: executeCrawlUseCase,
	}
}

// Execute handles POST /sites/:id/crawl
func (h *CrawlHandler) Execute(c *gin.Context) {
	siteID := c.Param("id")
	result, err := h.executeCrawlUseCase.Execute(c.Request.Context(), siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}


package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/crawler/internal/application/usecase/crawl"
)

// CrawlLogHandler handles HTTP requests for crawl logs
type CrawlLogHandler struct {
	getLogsUseCase      *crawl.GetLogsUseCase
	getLatestLogUseCase *crawl.GetLatestLogUseCase
}

// NewCrawlLogHandler creates a new CrawlLogHandler
func NewCrawlLogHandler(
	getLogsUseCase *crawl.GetLogsUseCase,
	getLatestLogUseCase *crawl.GetLatestLogUseCase,
) *CrawlLogHandler {
	return &CrawlLogHandler{
		getLogsUseCase:      getLogsUseCase,
		getLatestLogUseCase: getLatestLogUseCase,
	}
}

// GetLogs handles GET /sites/:id/logs
func (h *CrawlLogHandler) GetLogs(c *gin.Context) {
	siteID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	logs, err := h.getLogsUseCase.Execute(c.Request.Context(), siteID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// GetLatestLog handles GET /sites/:id/logs/latest
func (h *CrawlLogHandler) GetLatestLog(c *gin.Context) {
	siteID := c.Param("id")

	log, err := h.getLatestLogUseCase.Execute(c.Request.Context(), siteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": log})
}




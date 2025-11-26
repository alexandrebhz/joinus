package response

import (
	"github.com/gin-gonic/gin"
	"github.com/startup-job-board/backend/pkg/utils"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *utils.PaginationMeta `json:"meta,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, data interface{}, meta *utils.PaginationMeta) {
	c.JSON(200, SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}




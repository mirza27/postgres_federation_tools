package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetProgressRequest struct {
	Limit int `form:"limit" binding:"required"`
}

func (server *Server) GetProgress(c *gin.Context) {
	var req GetProgressRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	progressList, err := server.PivotDB.GetLastUpdatedQueueList(req.Limit)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to get progress list",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    progressList,
		"status":  "success",
		"message": "Successfully retrieved progress list",
	})
}

func (server *Server) GetProgressSummary(c *gin.Context) {

	summary, err := server.PivotDB.GetExecQueueSummary()
	if err != nil {
		c.JSON(400, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to get progress summary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    summary,
		"status":  "success",
		"message": "Successfully retrieved progress summary",
	})

}

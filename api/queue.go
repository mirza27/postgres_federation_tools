package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetExecutionQueueListRequest struct {
	Limit         int      `form:"limit" binding:"required"`
	Page          int      `form:"page" binding:"required"`
	SearchSQLText string   `form:"search.sql_text"`
	SearchSQLArg  string   `form:"search.sql_arg"`
	FilterStatus  []string `form:"filter.status"`
	FilterEntity  []string `form:"filter.entity"`
}

func (server *Server) GetExecutionQueueList(c *gin.Context) {
	var req GetExecutionQueueListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queueList, err := server.PivotDB.GetExecutionQueueList(req.Limit, req.Page, req.SearchSQLText, req.SearchSQLArg, req.FilterStatus, req.FilterEntity)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    queueList,
		"status":  "success",
		"message": "Successfully retrieved execution queue list",
	})
}

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

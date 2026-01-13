package api

import (
	"db_migrate_server/internal/debezium"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *Server) CheckAndCreateDebeziumConnector(c *gin.Context) {

	connectorName := server.Config.DebeziumConnectorName

	// check health
	dc, _ := debezium.NewDebeziumConnector(connectorName, server.Config, true)
	dcStatus, err := dc.Status()

	// if not found
	if dcStatus == "" {

		err := dc.Create()
		if err != nil {
			c.JSON(500, ErrorResponse{
				Error:   err.Error(),
				Status:  "error",
				Message: "Failed to create Debezium connector",
			})
			return
		}

		// delay get status
		time.Sleep(1 * time.Second)
		dcStatus, err = dc.Status()
	}

	// if connector exists, return success
	isConnectorExists := strings.Contains(dcStatus, connectorName)
	if !isConnectorExists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "Connector not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Connector checked/created successfully",
	})

}

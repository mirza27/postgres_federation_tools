package api

import (
	"db_migrate_server/internal/debezium"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *Server) CheckAndCreateDebeziumConnector(c *gin.Context) {

	connectorName := server.Config.DebeziumConnectorName

	// check health
	dc, _ := debezium.NewDebeziumConnector(connectorName, server.Config, true)
	dcStatus, err := dc.Status()

	// if contains 404, means connector not found
	if strings.Contains(dcStatus, "404") {
		fmt.Println("not found")
		// create connector if not exists
		err := dc.Create()
		if err != nil {
			c.JSON(500, ErrorResponse{
				Error:   err.Error(),
				Status:  "error",
				Message: "Failed to create Debezium connector",
			})
		}

		return
	}

	// if connector exists, return success
	if strings.Contains(dcStatus, connectorName) {
		c.JSON(http.StatusOK, DefaultResponse{
			Status:  "success",
			Message: "Connector checked/created successfully",
		})

		return
	}

	c.JSON(http.StatusOK, ErrorResponse{
		Status:  "error",
		Message: "Internal Server Error",
		Error:   err.Error(),
	})

}

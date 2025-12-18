package api

import (
	"context"
	"db_migrate_server/internal/pivot"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (server *Server) GetDatabaseConnection(c *gin.Context) {

	configList := map[string]any{
		"SOURCE_HOST":     "",
		"SOURCE_PORT":     int(0),
		"SOURCE_USER":     "",
		"SOURCE_PASSWORD": "",
		"SOURCE_DATABASE": "",
		"TARGET_HOST":     "",
		"TARGET_PORT":     int(0),
		"TARGET_USER":     "",
		"TARGET_PASSWORD": "",
		"TARGET_DATABASE": "",
		"TARGET_DSN":      "",
	}

	for key, _ := range configList {
		config, err := server.PivotDB.GetConfigurationByName(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   err.Error(),
				Status:  "error",
				Message: "Failed to get database configuration",
			})
			return
		}

		if config != nil {
			configList[key] = config.ConfigValue
		}
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Source database saved successfully",
		Data:    configList,
	})
}

type SaveDatabaseCredentialsRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	DbName   string `json:"db_name" binding:"required"`
}

func (server *Server) SaveSourceDatabase(c *gin.Context) {
	var req SaveDatabaseCredentialsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// save to pivot db
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "SOURCE_HOST",
		ConfigValue: req.Host,
	})
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "SOURCE_PORT",
		ConfigValue: fmt.Sprintf("%d", req.Port),
	})
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "SOURCE_USER",
		ConfigValue: req.User,
	})
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "SOURCE_PASSWORD",
		ConfigValue: req.Password,
	})
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "SOURCE_DATABASE",
		ConfigValue: req.DbName,
	})

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Source database saved successfully",
	})
}

func (server *Server) SaveTargetDatabase(c *gin.Context) {
	var req SaveDatabaseCredentialsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	targetDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		req.User, req.Password,
		req.Host, req.Port,
		req.DbName,
	)

	fmt.Println("Target DSN:", targetDSN)
	// check connection
	pool, err := pgxpool.New(context.Background(), targetDSN)
	pingErr := error(nil)
	if err == nil {
		pingErr = pool.Ping(context.Background())
	}
	if err != nil || pingErr != nil {
		e := err
		if e == nil {
			e = pingErr
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   e.Error(),
			Status:  "error",
			Message: "Failed to connect to target database",
		})
		return
	}
	pool.Close()

	// save to pivot db
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "TARGET_HOST",
		ConfigValue: req.Host,
	})
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "TARGET_PORT",
		ConfigValue: fmt.Sprintf("%d", req.Port),
	})
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "TARGET_USER",
		ConfigValue: req.User,
	})
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "TARGET_PASSWORD",
		ConfigValue: req.Password,
	})
	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "TARGET_DATABASE",
		ConfigValue: req.DbName,
	})

	server.PivotDB.UpdateConfigurationByName(&pivot.Configuration{
		ConfigKey:   "TARGET_DSN",
		ConfigValue: targetDSN,
	})

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Target database saved successfully",
	})
}

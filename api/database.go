package api

import (
	"context"
	"db_migrate_server/internal/pivot"
	"fmt"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type databaseConnectionResponse struct {
	SourceDatabase databaseCredentials `json:"source_database"`
	TargetDatabase databaseCredentials `json:"target_database"`
}
type databaseCredentials struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

func (server *Server) GetDatabaseConnection(c *gin.Context) {

	configList := map[string]any{
		"SOURCE_DATABASE_TYPE": "",
		"SOURCE_HOST":          "",
		"SOURCE_PORT":          int(0),
		"SOURCE_USER":          "",
		"SOURCE_PASSWORD":      "",
		"SOURCE_DATABASE":      "",
		"TARGET_DATABASE_TYPE": "",
		"TARGET_HOST":          "",
		"TARGET_PORT":          int(0),
		"TARGET_USER":          "",
		"TARGET_PASSWORD":      "",
		"TARGET_DATABASE":      "",
		"TARGET_DSN":           "",
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

	// database port
	source_port, _ := strconv.Atoi(configList["SOURCE_PORT"].(string))
	target_port, _ := strconv.Atoi(configList["TARGET_PORT"].(string))
	configList["SOURCE_PORT"] = source_port
	configList["TARGET_PORT"] = target_port

	// build response
	DatabaseResponse := databaseConnectionResponse{
		SourceDatabase: databaseCredentials{
			Type:     configList["SOURCE_DATABASE_TYPE"].(string),
			Host:     configList["SOURCE_HOST"].(string),
			Port:     configList["SOURCE_PORT"].(int),
			User:     configList["SOURCE_USER"].(string),
			Password: configList["SOURCE_PASSWORD"].(string),
			DbName:   configList["SOURCE_DATABASE"].(string),
		},
		TargetDatabase: databaseCredentials{
			Type:     configList["TARGET_DATABASE_TYPE"].(string),
			Host:     configList["TARGET_HOST"].(string),
			Port:     configList["TARGET_PORT"].(int),
			User:     configList["TARGET_USER"].(string),
			Password: configList["TARGET_PASSWORD"].(string),
			DbName:   configList["TARGET_DATABASE"].(string),
		},
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Source database saved successfully",
		Data:    DatabaseResponse,
	})
}

type SaveDatabaseCredentialsRequest struct {
	Type     string `json:"type" binding:"required"`
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

	targetDSN := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		req.Type, req.User, req.Password,
		req.Host, req.Port,
		req.DbName,
	)

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
		ConfigKey:   "SOURCE_DATABASE_TYPE",
		ConfigValue: req.Type,
	})
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

	targetDSN := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		req.Type, req.User, req.Password,
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
		ConfigKey:   "TARGET_DATABASE_TYPE",
		ConfigValue: req.Type,
	})
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

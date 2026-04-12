package api

import (
	"context"
	"db_migrate_server/internal/pivot"
	"fmt"
	"log"
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
		Message: "get database credentials successfully",
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

func (server *Server) GetSourceTables(c *gin.Context) {

	dbType := server.Config.TargetDatabaseType
	user := server.Config.TargetDatabaseUser
	password := server.Config.TargetDatabasePass
	host := server.Config.TargetDatabaseHost
	port := server.Config.TargetDatabasePort
	dbName := server.Config.TargetDatabaseName

	targetDSN := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		dbType, user, password, host, port, dbName,
	)

	log.Printf("Source DSN: %s", targetDSN)

	// check connection
	pool, err := pgxpool.New(c, targetDSN)
	pingErr := error(nil)
	if err == nil {
		pingErr = pool.Ping(c)
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

	q := `
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_type = 'BASE TABLE'
        ORDER BY table_name ASC;
    `
	tableNames, err := server.execDbQuery(pool, q)
	var result []map[string]interface{}

	for i := 0; i < len(tableNames.([]map[string]interface{})); i++ {

		tableName := tableNames.([]map[string]interface{})[i]["table_name"].(string)

		cq := fmt.Sprintf(`
			SELECT column_name, data_type 
			FROM information_schema.columns 
			WHERE table_schema = 'public' 
			AND table_name = '%s'
			ORDER BY ordinal_position;`,
			tableName)

		columns, _ := server.execDbQuery(pool, cq)

		result = append(result, map[string]interface{}{
			"table_name": tableName,
			"columns":    columns,
		})
	}

	// close connection
	pool.Close()

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to fetch table names",
		})
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Table names fetched successfully",
		Data:    result,
	})
}

func (server *Server) GetTargetTables(c *gin.Context) {

	dbType := server.Config.TargetDatabaseType
	user := server.Config.TargetDatabaseUser
	password := server.Config.TargetDatabasePass
	host := server.Config.TargetDatabaseHost
	port := server.Config.TargetDatabasePort
	dbName := server.Config.TargetDatabaseName

	targetDSN := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		dbType, user, password, host, port, dbName,
	)

	log.Printf("Target DSN: %s", targetDSN)

	// check connection
	pool, err := pgxpool.New(c, targetDSN)
	pingErr := error(nil)
	if err == nil {
		pingErr = pool.Ping(c)
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

	q := `
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_type = 'BASE TABLE'
        ORDER BY table_name ASC;
    `
	tableNames, err := server.execDbQuery(pool, q)
	var result []map[string]interface{}

	for i := 0; i < len(tableNames.([]map[string]interface{})); i++ {

		tableName := tableNames.([]map[string]interface{})[i]["table_name"].(string)

		cq := fmt.Sprintf(`
			SELECT column_name, data_type 
			FROM information_schema.columns 
			WHERE table_schema = 'public' 
			AND table_name = '%s'
			ORDER BY ordinal_position;`,
			tableName)

		columns, _ := server.execDbQuery(pool, cq)

		result = append(result, map[string]interface{}{
			"table_name": tableName,
			"columns":    columns,
		})
	}

	// close connection
	pool.Close()

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to fetch table names",
		})
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Table names fetched successfully",
		Data:    result,
	})
}

func (server *Server) GetDatabaseSchemaInfo(c *gin.Context) {
	// 1. Ambil semua konfigurasi dari PivotDB
	keys := []string{
		"SOURCE_DATABASE_TYPE", "SOURCE_USER", "SOURCE_PASSWORD", "SOURCE_HOST", "SOURCE_PORT", "SOURCE_DATABASE",
		"TARGET_DATABASE_TYPE", "TARGET_USER", "TARGET_PASSWORD", "TARGET_HOST", "TARGET_PORT", "TARGET_DATABASE",
	}

	configs := make(map[string]string)
	for _, key := range keys {
		cfg, err := server.PivotDB.GetConfigurationByName(key)
		if err != nil || cfg == nil {
			configs[key] = "" // fallback kosong
			continue
		}
		configs[key] = cfg.ConfigValue
	}

	// 2. Helper untuk mengambil schema dari DSN tertentu
	getSchema := func(dbType, user, pass, host, port, dbName string) (interface{}, error) {
		dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
			dbType, user, pass, host, port, dbName,
		)

		pool, err := pgxpool.New(c, dsn)
		if err != nil {
			return nil, err
		}
		defer pool.Close()

		if err := pool.Ping(c); err != nil {
			return nil, err
		}

		// Query ambil semua tabel
		q := `SELECT table_name FROM information_schema.tables 
              WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
              ORDER BY table_name ASC;`

		tableData, err := server.execDbQuery(pool, q)
		if err != nil {
			return nil, err
		}

		tables := tableData.([]map[string]interface{})
		var schemaResult []map[string]interface{}

		for _, t := range tables {
			tableName := t["table_name"].(string)

			// Query ambil kolom untuk tiap tabel
			cq := fmt.Sprintf(`
                SELECT column_name, data_type 
                FROM information_schema.columns 
                WHERE table_schema = 'public' AND table_name = '%s'
                ORDER BY ordinal_position;`, tableName)

			columns, _ := server.execDbQuery(pool, cq)
			schemaResult = append(schemaResult, map[string]interface{}{
				"table_name": tableName,
				"columns":    columns,
			})
		}
		return schemaResult, nil
	}

	// 3. Ambil data dari Source
	sourceSchema, err := getSchema(
		configs["SOURCE_DATABASE_TYPE"], configs["SOURCE_USER"], configs["SOURCE_PASSWORD"],
		configs["SOURCE_HOST"], configs["SOURCE_PORT"], configs["SOURCE_DATABASE"],
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect Source: " + err.Error()})
		return
	}

	// 4. Ambil data dari Target
	targetSchema, err := getSchema(
		configs["TARGET_DATABASE_TYPE"], configs["TARGET_USER"], configs["TARGET_PASSWORD"],
		configs["TARGET_HOST"], configs["TARGET_PORT"], configs["TARGET_DATABASE"],
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect Target: " + err.Error()})
		return
	}

	// 5. Kirim Response Gabungan
	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Database schema info fetched successfully",
		Data: map[string]interface{}{
			"source_schema": sourceSchema,
			"target_schema": targetSchema,
		},
	})
}

type GetColumnsRequest struct {
	TableName string `form:"table_name" binding:"required"`
}

func (server *Server) GetSourceColumns(c *gin.Context) {
	var req GetColumnsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request query",
		})
		return
	}

}

func (server *Server) GetTargetColumns(c *gin.Context) {
	var req GetColumnsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request query",
		})
		return
	}

}

func (server *Server) execDbQuery(pool *pgxpool.Pool, query string) (interface{}, error) {

	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	// create hashmap of column name
	columns := rows.FieldDescriptions()
	columnNames := make([]string, len(columns))
	for i, col := range columns {
		columnNames[i] = string(col.Name)
	}

	results := []map[string]interface{}{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columnNames {
			rowMap[colName] = values[i]
		}
		results = append(results, rowMap)
	}

	return results, nil
}

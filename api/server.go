package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"db_migrate_server/internal/config"
	"db_migrate_server/internal/pivot"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine  *gin.Engine
	Config  *config.Config
	PivotDB *pivot.Repo
}

func NewServer(pivotDb *pivot.Repo, config *config.Config) *Server {

	router := gin.Default()

	server := &Server{
		Engine:  router,
		Config:  config,
		PivotDB: pivotDb,
	}

	server.applyRoutes()

	return server
}

func (server *Server) applyRoutes() {

	server.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// mapping
	server.Engine.POST("/mapping", server.CreateMappingByName)
	server.Engine.PUT("/mapping", server.UpdateMappingByName)
	server.Engine.GET("/mapping", server.GetMappingByName)
	server.Engine.DELETE("/mapping", server.DeleteMappingByName)
	server.Engine.GET("/mapping/list", server.GetMappingList)

	// database
	server.Engine.GET("/database/", server.CheckDatabaseConnection)
	server.Engine.POST("/database/source", server.SaveSourceDatabase)
	server.Engine.POST("/database/target", server.SaveTargetDatabase)
	server.Engine.GET("/database/pivot", server.CheckPivotDatabases)

	// progress
	server.Engine.GET("/progress", server.GetProgress)
	server.Engine.GET("/progress/summary", server.GetProgressSummary)

	// worker
	server.Engine.POST("/worker/start", server.RunWorker)
	server.Engine.POST("/worker/stop", server.StopWorker)

}

func (server *Server) RunServer(ctx context.Context) error {

	addr := fmt.Sprintf(":%d", server.Config.RunApiPort)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: server.Engine,
	}
	log.Printf("Starting server on %s", addr)

	go func() {

		<-ctx.Done()
		_ = httpServer.Shutdown(context.Background())
	}()

	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

package main

import (
	"context"
	"db_migrate_server/api"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/pivot"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// get config
	config := config.Load()

	// get pivot db connection
	pivotDb, err := pivot.New(ctx, config.PivotDSN)
	if err != nil {
		log.Fatalf("failed to initialize pivot database: %v", err)
	}
	defer pivotDb.Close()

	log.Println("Pivot database initialized successfully")

	// run server
	apiServer := api.NewServer(pivotDb, config)
	if err := apiServer.RunServer(ctx); err != nil {
		log.Fatalf("failed to run API server: %v", err)
	}
}

package main

import (
	"context"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/executor"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	util.Info.Println("executor: signal context registered")
	defer util.Info.Println("executor: shutdown complete")

	util.Info.Println("executor: loading configuration")
	cfg := config.Load()
	util.Info.Println("executor: configuration loaded")

	util.Info.Println("executor: connecting pivot repository")
	pivotRepo, err := pivot.New(ctx, cfg.PivotDSN)
	if err != nil {
		panic(err)
	}
	defer pivotRepo.Close()
	util.Info.Println("executor: pivot repository connected")

	util.Info.Println("executor: ensuring pivot schema")
	if err := pivotRepo.EnsureSchema(ctx, pivot.DefaultSchemaSQL()); err != nil {
		panic(err)
	}
	util.Info.Println("executor: pivot schema ensured")

	util.Info.Printf("executor: creating executor batchMaxRows=%d intervalMs=%d", cfg.BatchMaxRows, cfg.BatchMaxInterval)
	exec, err := executor.New(ctx, pivotRepo, cfg.TargetDSN, cfg.BatchMaxRows, cfg.BatchMaxInterval)
	if err != nil {
		panic(err)
	}
	defer exec.Close()
	util.Info.Println("executor: executor ready")

	util.Info.Println("executor: starting run loop")
	exec.Run(ctx)
}

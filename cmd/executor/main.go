package main

import (
	"context"
	"db_migrate_server/internal/app"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/executor"
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
	config := config.Load()
	util.Info.Println("executor: configuration loaded")

	pivotRepo, err := app.InitPivot(ctx, config.PivotDSN)
	if err != nil {
		panic(err)
	}
	defer pivotRepo.Close()

	// apply base config
	config, err = app.ApplyBaseConfig(config, pivotRepo)
	if err != nil {
		panic(err)
	}

	// make executor
	util.Info.Printf("executor: creating executor batchMaxRows=%d intervalMs=%d", config.BatchMaxRows, config.BatchMaxInterval)
	exec, err := executor.New(ctx, pivotRepo, config.TargetDSN, config.BatchMaxRows, config.BatchMaxInterval)
	if err != nil {
		panic(err)
	}
	defer exec.Close()
	util.Info.Println("executor: executor ready")

	util.Info.Println("executor: starting run loop")
	exec.Run(ctx)
}

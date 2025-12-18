package main

import (
	"context"
	"db_migrate_server/internal/app"
	"db_migrate_server/internal/checker"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	util.Info.Println("checker: signal context registered")
	defer util.Info.Println("checker: shutdown complete")

	util.Info.Println("checker: loading configuration")
	config := config.Load()
	util.Info.Println("checker: configuration loaded")

	// check database pivot
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

	// load entity plan
	plan, err := app.LoadPlan(config)
	if err != nil {
		panic(err)
	}

	// make checker worker
	ch := checker.New(pivotRepo, plan, config.BatchMaxRows, config.BatchMaxInterval)
	util.Info.Println("checker: worker ready")

	util.Info.Println("checker: starting run loop")
	ch.Run(ctx)
}

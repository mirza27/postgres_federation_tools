package main

import (
	"context"
	"db_migrate_server/internal/app"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/joiner"
	"db_migrate_server/internal/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	util.Info.Println("joiner: signal context registered")
	defer util.Info.Println("joiner: shutdown complete")

	util.Info.Println("joiner: loading configuration")
	config := config.Load()
	util.Info.Println("joiner: configuration loaded")

	// make object connection to db pivot
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

	// make joiner worker
	w := joiner.New(pivotRepo, plan, config.BatchMaxRows, config.BatchMaxInterval)
	util.Info.Println("joiner: worker ready")

	util.Info.Println("joiner: starting run loop")
	w.Run(ctx)
}

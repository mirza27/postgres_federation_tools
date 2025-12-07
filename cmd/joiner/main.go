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
	cfg := config.Load()
	util.Info.Println("joiner: configuration loaded")

	plan, err := app.LoadPlan(cfg)
	if err != nil {
		panic(err)
	}

	pivotRepo, err := app.InitPivot(ctx, cfg.PivotDSN)
	if err != nil {
		panic(err)
	}
	defer pivotRepo.Close()

	w := joiner.New(pivotRepo, plan, cfg.BatchMaxRows, cfg.BatchMaxInterval)
	util.Info.Println("joiner: worker ready")

	util.Info.Println("joiner: starting run loop")
	w.Run(ctx)
}

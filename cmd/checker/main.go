package main

import (
	"context"
	"db_migrate_server/internal/checker"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/pivot"
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
	cfg := config.Load()
	util.Info.Println("checker: configuration loaded")

	util.Info.Println("checker: connecting pivot repository")
	pivotRepo, err := pivot.New(ctx, cfg.PivotDSN)
	if err != nil {
		panic(err)
	}
	defer pivotRepo.Close()
	util.Info.Println("checker: pivot repository connected")

	util.Info.Println("checker: ensuring pivot schema")
	if err := pivotRepo.EnsureSchema(ctx, pivot.DefaultSchemaSQL()); err != nil {
		panic(err)
	}
	util.Info.Println("checker: pivot schema ensured")

	ch := checker.New(pivotRepo, cfg.BatchMaxRows, cfg.BatchMaxInterval)
	util.Info.Println("checker: worker ready")

	util.Info.Println("checker: starting run loop")
	ch.Run(ctx)
}

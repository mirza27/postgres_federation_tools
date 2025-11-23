package main

import (
	"context"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pipeline"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	util.Info.Println("parser: signal context registered")
	defer util.Info.Println("parser: shutdown complete")

	util.Info.Println("parser: loading configuration")
	cfg := config.Load()
	util.Info.Println("parser: configuration loaded")

	util.Info.Println("parser: loading mapping root")
	root, err := mapping.Load(*cfg)
	if err != nil {
		panic(err)
	}
	util.Info.Printf("parser: mapping root loaded with %d entities", len(root.Entities))

	for i, e := range root.Entities {
		log.Printf("[%d] Entity: %s → target: %s", i+1, e.Entity, e.TargetTable)
	}

	util.Info.Println("parser: building planner")
	plan := mapping.NewPlanner(root)
	plan.Print()
	util.Info.Printf("parser: planner covers %d topics", len(plan.TopicList))

	util.Info.Printf("parser: creating kafka consumer group=%s", cfg.KafkaGroupID)
	consumer := kafka.NewConsumer([]string{cfg.KafkaBrokers}, cfg.KafkaGroupID, plan.TopicList)
	defer consumer.Close()
	util.Info.Println("parser: kafka consumer ready")

	util.Info.Println("parser: connecting pivot repository")
	pivotRepo, err := pivot.New(ctx, cfg.PivotDSN)
	if err != nil {
		panic(err)
	}
	defer pivotRepo.Close()
	util.Info.Println("parser: pivot repository connected")

	// util.Info.Println("parser: ensuring pivot schema")
	// if err := pivotRepo.EnsureSchema(ctx, pivot.DefaultSchemaSQL()); err != nil {
	// 	panic(err)
	// }
	// util.Info.Println("parser: pivot schema ensured")

	proc := pipeline.NewProcessor(plan, pivotRepo)
	util.Info.Println("parser: processor ready")

	util.Info.Println("parser: starting synchronous kafka loop")
	err = consumer.ListenSync(ctx, func(ev kafka.Event) error {
		if ev.Value == nil {
			util.Debug.Printf("parser: skip nil event topic=%s", ev.Topic)
			return nil
		}
		if err := proc.Handle(ctx, ev); err != nil {
			util.Error.Println("parser: process:", err)
		}
		return nil
	})
	if err != nil && ctx.Err() == nil {
		util.Error.Println("parser: kafka loop stopped:", err)
	}
}

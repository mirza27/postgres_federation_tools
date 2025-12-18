package main

import (
	"context"
	"db_migrate_server/internal/app"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/pipeline"
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
	config := config.Load()
	util.Info.Println("parser: configuration loaded")

	// check pivot db connection
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

	// check loaded plan
	for i, e := range plan.Entities {
		log.Printf("[%d] Entity: %s → target: %s", i+1, e.Entity, e.TargetTable)
	}
	plan.Print()
	util.Info.Printf("parser: planner covers %d topics", len(plan.TopicList))

	// make kafka consumer
	util.Info.Printf("parser: creating kafka consumer group=%s", config.KafkaGroupID)
	consumer := kafka.NewConsumer([]string{config.KafkaBrokers}, config.KafkaGroupID, plan.TopicList)
	defer consumer.Close()
	util.Info.Println("parser: kafka consumer ready")

	// make processor
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

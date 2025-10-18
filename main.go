package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"db_migrate_server/internal/config"
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pivot"
	// "github.com/yourorg/mapping-engine/internal/cache"
	// "github.com/yourorg/mapping-engine/internal/config"
	// "github.com/yourorg/mapping-engine/internal/executor"
	// "github.com/yourorg/mapping-engine/internal/kafka"
	// "github.com/yourorg/mapping-engine/internal/mapping"
	// "github.com/yourorg/mapping-engine/internal/pipeline"
	// "github.com/yourorg/mapping-engine/internal/pivot"
	// "github.com/yourorg/mapping-engine/internal/util"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// _, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// load env
	cfg := config.Load()

	// load config mapping
	root, err := mapping.Load(*cfg)
	if err != nil {
		panic(err)
	}

	for i, e := range root.Entities {
		log.Printf("[%d] Entity: %s → target: %s", i+1, e.Entity, e.TargetTable)
	}

	plan := mapping.NewPlanner(root)
	plan.Print()

	// set kafka connecttion
	consumer := kafka.NewConsumer([]string{cfg.KafkaBrokers}, cfg.KafkaGroupID, plan.TopicList)
	defer consumer.Close()

	// prepare db
	pivotRepo, err := pivot.New(ctx, cfg.PivotDSN)
	if err != nil {
		panic(err)
	}
	defer pivotRepo.Close()

	// apply schema in pivot db
	if err := pivotRepo.EnsureSchema(ctx, string(loadPivotSchema())); err != nil {
		panic(err)
	} else {
		log.Println("✅ Pivot schema ensured")
	}

	// jw := cache.NewJoinWait(cfg.JoinWaitTTL)
	// proc := pipeline.NewProcessor(plan, pivotRepo, jw)

	// exec, err := executor.New(ctx, pivotRepo, cfg.TargetDSN, cfg.BatchMaxRows, cfg.BatchMaxInterval)
	// if err != nil { panic(err) }
	// defer exec.Close()

	// // run executor background
	// go exec.Run(ctx)

	// // run consumer
	// ch := make(chan kafka.Event, 256)
	// go func() {
	// 	if err := consumer.Listen(ctx, ch); err != nil {
	// 		util.Error.Println("kafka listen:", err)
	// 		cancel()
	// 	}
	// }()

	// util.Info.Printf("Engine started. topics=%v", topicList)
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		util.Info.Println("shutting down")
	// 		return
	// 	case ev := <-ch:
	// 		if ev.Value == nil { continue }
	// 		if err := proc.Handle(ctx, ev); err != nil {
	// 			util.Error.Println("process:", err)
	// 		}
	// 	}
	// }
}

func loadPivotSchema() []byte {
	// pada implementasi nyata: baca dari file internal/pivot/schema.sql
	return []byte(`
create table if not exists _keymap_generic (
  map_name text not null,
  src_key text not null,
  tgt_key uuid not null,
  primary key (map_name, src_key),
  unique (map_name, tgt_key)
);
create table if not exists _exec_queue (
  id bigserial primary key,
  entity text not null,
  op text not null,
  sql_text text not null,
  sql_args jsonb,
  status text not null default 'pending',
  error text,
  created_at timestamptz default now()
);
create table if not exists _batch_log (
  id bigserial primary key,
  entity text not null,
  op text not null,
  key_values jsonb,
  payload jsonb,
  status text not null,
  error text,
  processed_at timestamptz default now(),
  batch_id text
);`)
}

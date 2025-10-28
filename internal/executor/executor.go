package executor

import (
	"context"
	"encoding/json"
	"time"

	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Executor struct {
	Pivot    *pivot.Repo
	Target   *pgxpool.Pool
	MaxRows  int
	Interval time.Duration
}

func New(ctx context.Context, pivotRepo *pivot.Repo, targetDSN string, maxRows, intervalMs int) (*Executor, error) {
	util.Info.Printf("executor: initializing target pool dsn=%s maxRows=%d interval=%dms", targetDSN, maxRows, intervalMs)
	tgt, err := pgxpool.New(ctx, targetDSN)
	if err != nil {
		return nil, err
	}
	return &Executor{
		Pivot: pivotRepo, Target: tgt,
		MaxRows: maxRows, Interval: time.Duration(intervalMs) * time.Millisecond,
	}, nil
}

func (e *Executor) Close() { e.Target.Close() }

func (e *Executor) Run(ctx context.Context) {
	ticker := time.NewTicker(e.Interval)
	defer ticker.Stop()
	util.Info.Printf("executor: run loop started interval=%s", e.Interval)

	for {
		select {
		case <-ctx.Done():
			util.Info.Println("executor: run loop canceled")
			return
		case <-ticker.C:
			e.tick(ctx)
		}
	}
}

func (e *Executor) tick(ctx context.Context) {
	items, err := e.Pivot.FetchPending(ctx, e.MaxRows)
	if err != nil {
		util.Error.Println("FetchPending:", err)
		return
	}
	if len(items) == 0 {
		util.Debug.Println("executor: no pending items")
		return
	}
	util.Info.Printf("executor: processing %d pending items", len(items))
	tx, err := e.Target.Begin(ctx)
	if err != nil {
		util.Error.Println("begin tx:", err)
		return
	}
	defer tx.Rollback(ctx)

	for _, it := range items {
		util.Debug.Printf("executor: handling queue id=%d entity=%s op=%s", it.ID, it.Entity, it.Op)
		var args []interface{}
		if err := json.Unmarshal(it.ArgsJSON, &args); err != nil {
			util.Error.Println("args unmarshal:", err)
			_ = e.Pivot.MarkError(ctx, it.ID, err.Error())
			continue
		}
		_, err := tx.Exec(ctx, it.SQL, args...)
		if err != nil {
			util.Error.Printf("exec id=%d err=%v sql=%s", it.ID, err, it.SQL)
			_ = e.Pivot.MarkError(ctx, it.ID, err.Error())
			continue
		}
		if err := e.Pivot.MarkDone(ctx, it.ID); err != nil {
			util.Warn.Println("mark done:", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		util.Error.Println("commit:", err)
	} else {
		util.Info.Printf("executor: committed batch items=%d", len(items))
	}
}

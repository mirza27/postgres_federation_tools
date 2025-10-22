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

	for {
		select {
		case <-ctx.Done():
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
		return
	}

	tx, err := e.Target.Begin(ctx)
	if err != nil {
		util.Error.Println("begin tx:", err)
		return
	}
	defer tx.Rollback(ctx)

	for _, it := range items {
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
	}
}

package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Executor struct {
	Pivot    *pivot.Repo
	Target   *pgxpool.Pool
	MaxRows  int
	Interval time.Duration
	WorkerID string
}

func New(ctx context.Context, pivotRepo *pivot.Repo, targetDSN string, maxRows, intervalMs int) (*Executor, error) {
	util.Info.Printf("executor: initializing target pool dsn=%s maxRows=%d interval=%dms", targetDSN, maxRows, intervalMs)
	tgt, err := pgxpool.New(ctx, targetDSN)
	if err != nil {
		return nil, err
	}
	return &Executor{
		Pivot:    pivotRepo,
		Target:   tgt,
		MaxRows:  maxRows,
		Interval: time.Duration(intervalMs) * time.Millisecond,
		WorkerID: fmt.Sprintf("executor-%s", uuid.NewString()),
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
	items, err := e.Pivot.FetchReady(ctx, e.MaxRows)
	if err != nil {
		util.Error.Println("FetchReady:", err)
		return
	}
	if len(items) == 0 {
		util.Debug.Println("executor: no ready items")
		return
	}
	util.Info.Printf("executor: processing %d ready items", len(items))
	tx, err := e.Target.Begin(ctx)
	if err != nil {
		util.Error.Println("begin tx:", err)
		return
	}
	defer tx.Rollback(ctx)

	for _, it := range items {
		util.Debug.Printf("executor: handling queue id=%d entity=%s op=%s needKeymap=%t", it.ID, it.Entity, it.Op, it.NeedKeymap)

		if err := e.Pivot.MarkExecuting(ctx, it.ID, e.WorkerID); err != nil {
			util.Warn.Printf("executor: mark executing failed id=%d err=%v", it.ID, err)
			continue
		}

		var args []interface{}
		if err := json.Unmarshal(it.ArgsJSON, &args); err != nil {
			util.Error.Printf("executor: args unmarshal failed id=%d err=%v", it.ID, err)
			_ = e.Pivot.MarkError(ctx, it.ID, err.Error())
			if it.KeymapID != nil {
				_ = e.Pivot.MarkKeymapError(ctx, *it.KeymapID, err.Error())
			}
			continue
		}

		if err := e.executeItem(ctx, tx, it, args); err != nil {
			util.Error.Printf("executor: execute failed id=%d err=%v", it.ID, err)
			_ = e.Pivot.MarkError(ctx, it.ID, err.Error())
			if it.KeymapID != nil {
				_ = e.Pivot.MarkKeymapError(ctx, *it.KeymapID, err.Error())
			}
			continue
		}

		if err := e.Pivot.MarkDone(ctx, it.ID); err != nil {
			util.Warn.Printf("executor: mark done failed id=%d err=%v", it.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		util.Error.Println("commit:", err)
	} else {
		util.Info.Printf("executor: committed batch items=%d", len(items))
	}
}

func (e *Executor) executeItem(ctx context.Context, tx pgx.Tx, it pivot.Row, args []interface{}) error {
	if it.NeedKeymap && len(it.Returning) == 0 {
		util.Warn.Printf("executor: queue id=%d marked needKeymap but returning columns missing", it.ID)
	}

	if it.NeedKeymap && len(it.Returning) > 0 {
		dest := make([]interface{}, len(it.Returning))
		scanTargets := make([]interface{}, len(it.Returning))
		for i := range dest {
			scanTargets[i] = &dest[i]
		}
		if err := tx.QueryRow(ctx, it.SQL, args...).Scan(scanTargets...); err != nil {
			return err
		}
		if len(dest) > 0 && it.KeymapID != nil {
			key := fmt.Sprint(dest[0])
			if err := e.Pivot.FulfillKeymap(ctx, *it.KeymapID, key); err != nil {
				return err
			}
		} else if len(dest) > 0 && it.KeymapID == nil {
			util.Warn.Printf("executor: queue id=%d returning key but keymap_id nil", it.ID)
		}
		return nil
	}

	_, err := tx.Exec(ctx, it.SQL, args...)
	return err
}

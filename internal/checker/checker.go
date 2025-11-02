package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"

	"github.com/google/uuid"
)

type Checker struct {
	Pivot    *pivot.Repo
	MaxRows  int
	Interval time.Duration
	WorkerID string
}

func New(pivotRepo *pivot.Repo, maxRows, intervalMs int) *Checker {
	return &Checker{
		Pivot:    pivotRepo,
		MaxRows:  maxRows,
		Interval: time.Duration(intervalMs) * time.Millisecond,
		WorkerID: fmt.Sprintf("checker-%s", uuid.NewString()),
	}
}

func (c *Checker) Run(ctx context.Context) {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	util.Info.Printf("checker: run loop started interval=%s", c.Interval)

	for {
		select {
		case <-ctx.Done():
			util.Info.Println("checker: run loop canceled")
			return
		case <-ticker.C:
			c.tick(ctx)
		}
	}
}

func (c *Checker) tick(ctx context.Context) {
	items, err := c.Pivot.FetchForChecker(ctx, c.MaxRows)
	if err != nil {
		util.Error.Println("checker: fetch:", err)
		return
	}
	if len(items) == 0 {
		util.Debug.Println("checker: no pending items")
		return
	}
	util.Info.Printf("checker: processing %d items", len(items))

	for _, item := range items {
		util.Debug.Printf("checker: handling queue id=%d entity=%s needKeymap=%t", item.ID, item.Entity, item.NeedKeymap)
		if !item.NeedKeymap {
			if err := c.Pivot.MarkReady(ctx, item.ID); err != nil {
				util.Warn.Printf("checker: mark ready failed id=%d err=%v", item.ID, err)
			}
			continue
		}

		var req pivot.KeymapRequest
		if len(item.KeymapPayload) > 0 {
			if err := json.Unmarshal(item.KeymapPayload, &req); err != nil {
				util.Error.Printf("checker: unmarshall keymap payload failed id=%d err=%v", item.ID, err)
				_ = c.Pivot.MarkError(ctx, item.ID, err.Error())
				continue
			}
		}
		if req.MapName == "" {
			req.MapName = item.Entity
		}
		if req.SrcKey == "" {
			util.Warn.Printf("checker: queue id=%d missing src_key for keymap request", item.ID)
			_ = c.Pivot.MarkError(ctx, item.ID, "missing src_key for keymap request")
			continue
		}

		keymapID, err := c.Pivot.EnsureKeymapRequest(ctx, item.QueueID, req)
		if err != nil {
			util.Error.Printf("checker: ensure keymap failed id=%d err=%v", item.ID, err)
			_ = c.Pivot.MarkError(ctx, item.ID, err.Error())
			continue
		}

		if err := c.Pivot.AttachKeymap(ctx, item.ID, keymapID, "ready"); err != nil {
			util.Warn.Printf("checker: attach keymap failed id=%d err=%v", item.ID, err)
			continue
		}
	}
}

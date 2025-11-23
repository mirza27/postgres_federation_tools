package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pipeline"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"

	"github.com/google/uuid"
)

type Checker struct {
	Pivot    *pivot.Repo
	Plan     *mapping.Planner
	Proc     *pipeline.Processor
	MaxRows  int
	Interval time.Duration
	WorkerID string
}

func New(pivotRepo *pivot.Repo, plan *mapping.Planner, maxRows, intervalMs int) *Checker {
	return &Checker{
		Pivot:    pivotRepo,
		Plan:     plan,
		Proc:     pipeline.NewProcessor(plan, pivotRepo),
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
		util.Debug.Printf("checker: handling queue id=%d entity=%s needKeymap=%t needJoin=%t", item.ID, item.Entity, item.NeedKeymap, item.NeedJoin)

		if item.NeedJoin {
			if err := c.handleJoin(ctx, item); err != nil {
				util.Warn.Printf("checker: handle join failed id=%d err=%v", item.ID, err)
			}
			continue
		}

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

func (c *Checker) handleJoin(ctx context.Context, item pivot.CheckerItem) error {
	ent, ok := c.findEntity(item.Entity)
	if !ok {
		return fmt.Errorf("entity %s not found in planner", item.Entity)
	}
	expectedTopics := mapping.ExpectedTopics(ent)

	// fact payload dari queue
	var factPayload kafka.DebeziumValue
	if len(item.JoinPayload) > 0 {
		if err := json.Unmarshal(item.JoinPayload, &factPayload); err != nil {
			return err
		}
	}

	// dimTopics: semua topic kecuali topic faktual yang ada di queue
	var dimTopics []string
	for _, t := range expectedTopics {
		if t != item.JoinTopic {
			dimTopics = append(dimTopics, t)
		}
	}
	if len(dimTopics) == 0 {
		// tidak ada dimensi, langsung enqueue
		return c.Proc.HandleJoinReady(ctx, ent, item.Op, map[string]*kafka.DebeziumValue{
			item.JoinTopic: &factPayload,
		})
	}

	fragments, complete, err := c.Pivot.FetchJoinFragments(ctx, item.Entity, item.JoinKey, dimTopics)
	if err != nil {
		return err
	}
	if !complete {
		util.Debug.Printf("checker: join incomplete entity=%s joinKey=%s", item.Entity, item.JoinKey)
		return nil
	}

	payloads := map[string]*kafka.DebeziumValue{
		item.JoinTopic: &factPayload,
	}
	for t, p := range fragments {
		payloads[t] = p
	}

	if err := c.Proc.HandleJoinReady(ctx, ent, item.Op, payloads); err != nil {
		return err
	}
	// tandai placeholder selesai
	return c.Pivot.MarkDone(ctx, item.ID)
}

func (c *Checker) findEntity(name string) (mapping.Entity, bool) {
	for _, e := range c.Plan.Entities {
		if e.Entity == name {
			return e, true
		}
	}
	return mapping.Entity{}, false
}

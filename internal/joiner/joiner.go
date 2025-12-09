package joiner

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pipeline"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"

	"github.com/google/uuid"
)

// Worker joiner mengambil pekerjaan join pending, menunggu fragment lengkap,
// lalu menulis baris eksekusi nyata ke _exec_queue.
type Worker struct {
	Pivot    PivotStore
	Plan     *mapping.Planner
	Proc     *pipeline.Processor
	MaxRows  int
	Interval time.Duration
	WorkerID string
}

type PivotStore interface {
	pipeline.PivotStore
	FetchNeedJoin(ctx context.Context, limit int) ([]pivot.NeedJoinItem, error)
	FetchJoinFragments(ctx context.Context, entity, joinKey string, topics []string) (map[string]*kafka.DebeziumValue, bool, error)
	MarkNeedJoinError(ctx context.Context, id int64, msg string) error
	MarkNeedJoinDone(ctx context.Context, id int64) error
}

func New(pivotRepo PivotStore, plan *mapping.Planner, maxRows, intervalMs int) *Worker {
	return &Worker{
		Pivot:    pivotRepo,
		Plan:     plan,
		Proc:     pipeline.NewProcessor(plan, pivotRepo),
		MaxRows:  maxRows,
		Interval: time.Duration(intervalMs) * time.Millisecond,
		WorkerID: fmt.Sprintf("joiner-%s", uuid.NewString()),
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	util.Info.Printf("joiner: run loop started interval=%s", w.Interval)

	for {
		select {
		case <-ctx.Done():
			util.Info.Println("joiner: run loop canceled")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	items, err := w.Pivot.FetchNeedJoin(ctx, w.MaxRows)
	if err != nil {
		util.Error.Println("joiner: fetch:", err)
		return
	}
	if len(items) == 0 {
		util.Debug.Println("joiner: no pending items")
		return
	}
	util.Info.Printf("joiner: processing %d items", len(items))

	for _, item := range items {
		if err := w.handle(ctx, item); err != nil {
			util.Warn.Printf("joiner: handle failed id=%d err=%v", item.ID, err)
			_ = w.Pivot.MarkNeedJoinError(ctx, item.ID, err.Error())
		}
	}
}

func (w *Worker) handle(ctx context.Context, item pivot.NeedJoinItem) error {

	ent, ok := w.findEntity(item.Entity)
	if !ok {
		return fmt.Errorf("entity %s not found in planner", item.Entity)
	}
	expectedTopics := mapping.ExpectedTopics(ent)

	// fact payload dari need_join
	var factPayload kafka.DebeziumValue
	if len(item.JoinPayload) > 0 {
		if err := json.Unmarshal(item.JoinPayload, &factPayload); err != nil {
			return err
		}
	}

	// dimTopics: semua topic kecuali topic faktual yang ada di need_join
	var dimTopics []string
	for _, t := range expectedTopics {
		if t != item.JoinTopic {
			dimTopics = append(dimTopics, t)
		}
	}
	if len(dimTopics) == 0 {
		// tidak ada dimensi, langsung enqueue
		if err := w.Proc.HandleJoinReadyWithQueue(ctx, ent, item.Op, item.QueueID, map[string]*kafka.DebeziumValue{
			item.JoinTopic: &factPayload,
		}); err != nil {
			return err
		}
		return w.Pivot.MarkNeedJoinDone(ctx, item.ID)
	}

	fragments, complete, err := w.Pivot.FetchJoinFragments(ctx, item.Entity, item.JoinKey, dimTopics)
	if err != nil {
		return err
	}

	// jika join_key komposit (dipisah '|'), coba ambil fragmen per bagian dan gabungkan
	allFragments := map[string]*kafka.DebeziumValue{}
	if len(item.JoinFields) > 0 {
		var jf map[string]string
		if err := json.Unmarshal(item.JoinFields, &jf); err == nil {
			// gunakan keys dari join_fields
			for k := range jf {
				partsFrag, _, ferr := w.Pivot.FetchJoinFragments(ctx, item.Entity, k, dimTopics)
				if ferr != nil {
					return ferr
				}
				for t, p := range partsFrag {
					allFragments[t] = p
				}
			}
		}
	}
	// fallback: split join_key dengan '|'
	if len(allFragments) == 0 && strings.Contains(item.JoinKey, "|") {
		for _, k := range strings.Split(item.JoinKey, "|") {
			k = strings.TrimSpace(k)
			if k == "" {
				continue
			}
			partsFrag, _, ferr := w.Pivot.FetchJoinFragments(ctx, item.Entity, k, dimTopics)
			if ferr != nil {
				return ferr
			}
			for t, p := range partsFrag {
				allFragments[t] = p
			}
		}
	}
	// jika tidak ada hasil gabungan, pakai hasil fetch awal
	if len(allFragments) == 0 {
		for t, p := range fragments {
			allFragments[t] = p
		}
	}

	// cek kelengkapan manual
	complete = true
	for _, t := range dimTopics {
		if _, ok := allFragments[t]; !ok {
			complete = false
			break
		}
	}
	if !complete {
		util.Debug.Printf("joiner: join incomplete entity=%s joinKey=%s", item.Entity, item.JoinKey)
		return nil
	}

	payloads := map[string]*kafka.DebeziumValue{
		item.JoinTopic: &factPayload,
	}
	for t, p := range allFragments {
		payloads[t] = p
	}

	if err := w.Proc.HandleJoinReadyWithQueue(ctx, ent, item.Op, item.QueueID, payloads); err != nil {
		return err
	}
	return w.Pivot.MarkNeedJoinDone(ctx, item.ID)
}

func (w *Worker) findEntity(name string) (mapping.Entity, bool) {
	// util.Debug.Printf("name %s in entity %s", name, w.Plan.Entities)

	for _, e := range w.Plan.Entities {
		if e.Entity == name {
			return e, true
		}
	}
	return mapping.Entity{}, false
}

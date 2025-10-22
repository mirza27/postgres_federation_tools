package pipeline

import (
	"context"
	"db_migrate_server/internal/cache"
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/sqlbuilder"
	"db_migrate_server/internal/util"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

type Processor struct {
	Plan  *mapping.Planner
	Pivot *pivot.Repo
	Join  *cache.JoinWait
}

func NewProcessor(plan *mapping.Planner, pivotRepo *pivot.Repo, jw *cache.JoinWait) *Processor {
	return &Processor{Plan: plan, Pivot: pivotRepo, Join: jw}
}

func (p *Processor) Handle(ctx context.Context, ev kafka.Event) error {

	op, _ := kafka.ExtractOpAndSource(&ev)
	entities := p.Plan.TopicToEntities[ev.Topic]

	if len(entities) == 0 {
		util.Debug.Printf("no entity matched topic=%s", ev.Topic)
		return nil
	}

	for _, ent := range entities {
		// 1) join-wait jika perlu lebih dari 1 topic
		allTopics := mapping.ExpectedTopics(ent)

		if len(allTopics) > 1 {

			joinKey := p.deriveJoinKey(ent, ev.Value)
			p.Join.Add(joinKey, ev.Topic, ev.Value)

			if vals, ok := p.Join.GetSet(joinKey, allTopics); ok {
				// compose merged payload di memori
				merged := p.mergePayload(vals)
				if err := p.planAndEnqueue(ctx, ent, op, merged); err != nil {
					return err
				}
			} else {
				continue // tunggu event lain
			}

		} else {
			if err := p.planAndEnqueue(ctx, ent, op, map[string][]byte{ev.Topic: ev.Value}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) deriveJoinKey(ent mapping.Entity, value []byte) string {

	// Ambil key.source[0] misalnya "u.id" → field "id" di after
	// Disederhanakan: baca from after.id (Debezium)

	v := gjson.ParseBytes(value)
	after := v.Get("after")
	if !after.Exists() {
		after = v.Get("payload.after")
	}

	key := ""
	if len(ent.Key.Source) > 0 {
		// ambil path terakhir (kolom)
		src := ent.Key.Source[0] // "u.id"
		parts := strings.Split(src, ".")
		col := parts[len(parts)-1]
		key = after.Get(col).String()
	}

	return key

}

// mergePayload: gabungkan payload dari beberapa topic (sederhana: map[topic]rawJSON)
func (p *Processor) mergePayload(topicBytes map[string][]byte) map[string][]byte { return topicBytes }

func (p *Processor) planAndEnqueue(ctx context.Context, ent mapping.Entity, op string, payload map[string][]byte) error {
	// 2) Resolve key: natural/shared/surrogate (sederhana)
	keyVal, err := p.resolveKey(ctx, ent, payload)
	if err != nil {
		return err
	}

	// 3) Bangun kolom target
	cols, vals, sets, where, conflictCols, err := p.buildColumns(ent, keyVal, payload)
	if err != nil {
		return err
	}

	// 4) Tentukan writeMode & buat SQL
	var stmt sqlbuilder.Stmt
	switch opRoute(op, ent) {
	case "insert", "upsert":
		if ent.Routing.OnCreate.WriteMode == "upsert" || ent.Routing.OnUpdate.WriteMode == "upsert" || ent.Routing.OnSnapshot.WriteMode == "upsert" {
			updates := make([]string, 0, len(cols))
			for i, c := range cols {
				updates = append(updates, fmt.Sprintf(`%s=EXCLUDED.%s`, c, c))
				_ = i
			}
			stmt = sqlbuilder.Upsert(ent.TargetTable, cols, vals, conflictCols, updates)
		} else {
			stmt = sqlbuilder.Insert(ent.TargetTable, cols, vals)
		}
	case "update":
		stmt = sqlbuilder.Update(ent.TargetTable, sets, where, vals)
	case "delete":
		stmt = sqlbuilder.Delete(ent.TargetTable, where, vals)
	default:
		return fmt.Errorf("unsupported route")
	}

	// 5) Enqueue ke pivot._exec_queue
	return p.Pivot.Enqueue(ctx, pivot.ExecItem{
		Entity: ent.Entity, Op: op, SQL: stmt.SQL, Args: stmt.Args,
	})
}

func opRoute(op string, ent mapping.Entity) string {
	switch op {
	case "c":
		return ent.Routing.OnCreate.WriteMode
	case "u":
		return ent.Routing.OnUpdate.WriteMode
	case "d":
		return ent.Routing.OnDelete.WriteMode
	case "r":
		return ent.Routing.OnSnapshot.WriteMode
	default:
		return "upsert"
	}
}

func (p *Processor) resolveKey(ctx context.Context, ent mapping.Entity, payload map[string][]byte) (string, error) {
	switch ent.Key.Strategy {
	case "natural":
		// ambil dari after kolom pertama (sederhana)
		return p.deriveJoinKey(ent, anyPayload(payload)), nil
	case "shared_key", "surrogate":
		if ent.Key.Resolver == nil || ent.Key.Resolver.Table == "" {
			return "", fmt.Errorf("missing resolver for %s", ent.Entity)
		}
		mapName := ent.Key.Resolver.Table // atau pakai nama entitas (bebas)
		srcKey := p.deriveJoinKey(ent, anyPayload(payload))
		// lookup
		if tgt, ok, _ := p.Pivot.LookupKey(ctx, mapName, srcKey); ok {
			return tgt, nil
		}
		// generate (uuid v7 di sisi Go; di contoh sederhana ini pakai pg gen_random_uuid() boleh juga)
		// gen := GenUUIDv7()
		gen := uuid.New().String()
		return p.Pivot.PutKeyIfAbsent(ctx, mapName, srcKey, gen)
	default:
		return "", fmt.Errorf("unknown key strategy")
	}
}

func anyPayload(m map[string][]byte) []byte {
	for _, v := range m {
		return v
	}
	return nil
}

func (p *Processor) buildColumns(ent mapping.Entity, key string, payload map[string][]byte) (
	cols []string, vals []interface{}, sets []string, where []string, conflict []string, err error) {

	// Disederhanakan: hanya handle "$key" dan "from: <col>" dari after.* ; "expr: now()" diisi di SQL (skipped) atau di Go (time.Now())
	// Kamu bisa perluasan: cast, resolver lookup, flatten/aggregate (untuk contoh ringkas ini belum penuh).
	for colName, spec := range ent.Columns {
		cols = append(cols, colName)
		switch {
		case spec.From == "$key":
			vals = append(vals, key)
		case spec.Expr == "now()":
			vals = append(vals, time.Now().UTC())
		case spec.From != "":
			// cari dari payload topic mana pun: ambil field "after.<col>"
			vals = append(vals, firstAfterField(payload, spec.From))
		default:
			vals = append(vals, spec.Default)
		}
	}

	// sets untuk UPDATE (col=$i), where dari matchKey atau primary (key)
	for i, c := range cols {
		sets = append(sets, fmt.Sprintf("%s=$%d", c, i+1))
	}

	if len(ent.Routing.OnCreate.ConflictKey) > 0 {
		conflict = ent.Routing.OnCreate.ConflictKey
	} else if len(ent.Routing.OnUpdate.ConflictKey) > 0 {
		conflict = ent.Routing.OnUpdate.ConflictKey
	} else if len(ent.Routing.OnSnapshot.ConflictKey) > 0 {
		conflict = ent.Routing.OnSnapshot.ConflictKey
	}

	// where sederhana: if MatchKey ada → gunakan itu, else pakai key kolom pertama di conflict
	mk := ent.Routing.OnDelete.MatchKey
	if len(mk) == 0 && len(conflict) > 0 {
		mk = conflict
	}
	for idx, k := range mk {
		where = append(where, fmt.Sprintf("%s=$%d", k, idx+1))
	}

	return
}

func firstAfterField(payload map[string][]byte, path string) any {
	// path bisa "u.email" / "o.id" → ambil segmen terakhir
	parts := strings.Split(path, ".")
	col := parts[len(parts)-1]
	for _, raw := range payload {
		v := gjson.ParseBytes(raw)
		a := v.Get("after")
		if !a.Exists() {
			a = v.Get("payload.after")
		}
		if a.Exists() {
			x := a.Get(col)
			if x.Exists() {
				return x.Value()
			}
		}
	}
	return nil
}

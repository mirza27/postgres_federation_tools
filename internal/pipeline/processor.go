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

	op := ev.Op
	if op == "" && ev.Value != nil {
		op = ev.Value.Op
	}
	entities := p.Plan.TopicToEntities[ev.Topic]

	util.Debug.Printf("processor: Handle topic=%s op=%s matchedEntities=%d", ev.Topic, op, len(entities))

	if len(entities) == 0 {
		util.Debug.Printf("no entity matched topic=%s", ev.Topic)
		return nil
	}

	for _, ent := range entities {

		util.Debug.Printf("processor: handling entity=%s op=%s", ent.Entity, op)

		// join wait if entity need more than 1 topic
		allTopics := mapping.ExpectedTopics(ent)
		util.Debug.Printf("processor: expected topics=%v for entity=%s", allTopics, ent.Entity)

		if len(allTopics) > 1 {

			joinKey := p.deriveJoinKey(ent, ev.Value)
			util.Debug.Printf("processor: join key=%s topic=%s", joinKey, ev.Topic)
			p.Join.Add(joinKey, ev.Topic, ev.Value)

			if vals, ok := p.Join.GetSet(joinKey, allTopics); ok {
				// compose merged payload di memori
				util.Debug.Printf("processor: join satisfied entity=%s topics=%v", ent.Entity, allTopics)
				merged := p.mergePayload(vals)
				if err := p.planAndEnqueue(ctx, ent, op, merged); err != nil {
					return err
				}
			} else {
				util.Debug.Printf("processor: join incomplete entity=%s waiting other topics", ent.Entity)
				continue // tunggu event lain
			}

			// if entity need only one topic
		} else {
			util.Debug.Printf("processor: single topic path entity=%s topic=%s", ent.Entity, ev.Topic)

			if err := p.planAndEnqueue(ctx, ent, op, map[string]*kafka.DebeziumValue{ev.Topic: ev.Value}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) deriveJoinKey(ent mapping.Entity, value *kafka.DebeziumValue) string {

	// Ambil key.source[0] misalnya "u.id" → field "id" di after
	util.Debug.Printf("processor: deriveJoinKey entity=%s", ent.Entity)

	if value == nil || len(ent.Key.Source) == 0 {
		util.Debug.Printf("processor: deriveJoinKey missing value or key source entity=%s", ent.Entity)
		return ""
	}

	src := ent.Key.Source[0] // "u.id"
	parts := strings.Split(src, ".")
	col := parts[len(parts)-1]

	key := stringFromRow(value.Payload.After, col)
	if key == "" {
		key = stringFromRow(value.Payload.Before, col)
	}

	util.Debug.Printf("processor: derived join key=%s entity=%s", key, ent.Entity)
	return key
}

// mergePayload: gabungkan payload dari beberapa topic (sederhana: map[topic]rawJSON)
func (p *Processor) mergePayload(topicBytes map[string]*kafka.DebeziumValue) map[string]*kafka.DebeziumValue {
	return topicBytes
}

func (p *Processor) planAndEnqueue(ctx context.Context, ent mapping.Entity, op string, payload map[string]*kafka.DebeziumValue) error {

	util.Debug.Printf("processor: planAndEnqueue entity=%s op=%s payloadTopics=%d", ent.Entity, op, len(payload))

	// check foreign type
	keyVal, err := p.resolveKey(ctx, ent, payload)
	if err != nil {
		return err
	}
	util.Debug.Printf("processor: resolved key=%s entity=%s", keyVal, ent.Entity)

	// 3) Bangun kolom target
	cols, vals, sets, where, conflictCols, err := p.buildColumns(ent, keyVal, payload)
	if err != nil {
		return err
	}
	util.Debug.Printf("processor: built columns entity=%s cols=%v conflict=%v", ent.Entity, cols, conflictCols)

	// 4) Tentukan writeMode & buat SQL
	var stmt sqlbuilder.Stmt

	switch opRoute(op, ent) { // check event op type
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
		util.Debug.Printf("processor: delete route entity=%s where=%v", ent.Entity, where)

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

// mendefinisikan mapping foreign berdasarkan 'strategy'
func (p *Processor) resolveKey(ctx context.Context, ent mapping.Entity, payload map[string]*kafka.DebeziumValue) (string, error) {
	util.Debug.Printf("processor: resolveKey entity=%s strategy=%s", ent.Entity, ent.Key.Strategy)

	switch ent.Key.Strategy {

	case "natural":

		evtPayload := anyPayload(payload)

		key := p.deriveJoinKey(ent, evtPayload)

		return key, nil

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

func anyPayload(m map[string]*kafka.DebeziumValue) *kafka.DebeziumValue {
	for _, v := range m {
		util.Debug.Printf("processor: anyPayload returning payload topic")

		return v
	}
	util.Debug.Printf("processor: anyPayload no payload available")
	return nil
}

func stringFromRow(row map[string]any, key string) string {
	if row == nil {
		return ""
	}
	if val, ok := row[key]; ok && val != nil {
		switch v := val.(type) {
		case string:
			return v
		case fmt.Stringer:
			return v.String()
		default:
			return fmt.Sprint(v)
		}
	}
	return ""
}

func (p *Processor) buildColumns(ent mapping.Entity, key string, payload map[string]*kafka.DebeziumValue) (
	cols []string, vals []interface{}, sets []string, where []string, conflict []string, err error) {

	// Disederhanakan: hanya handle "$key" dan "from: <col>" dari after.* ; "expr: now()" diisi di SQL (skipped) atau di Go (time.Now())
	// Kamu bisa perluasan: cast, resolver lookup, flatten/aggregate (untuk contoh ringkas ini belum penuh).

	fmt.Printf("build column ent %v\n", ent)
	fmt.Printf("build column ent columns %v\n", ent.Columns)
	fmt.Printf("build column payload %v\n", payload)

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

func firstAfterField(payload map[string]*kafka.DebeziumValue, path string) any {
	// path bisa "u.email" / "o.id" → ambil segmen terakhir
	util.Debug.Printf("processor: firstAfterField path=%s", path)
	parts := strings.Split(path, ".")
	col := parts[len(parts)-1]
	for _, evt := range payload {
		if evt == nil {
			continue
		}
		if evt.Payload.After != nil {
			if val, ok := evt.Payload.After[col]; ok {
				return val
			}
		}
		if evt.Payload.Before != nil {
			if val, ok := evt.Payload.Before[col]; ok {
				return val
			}
		}
	}
	return nil
}

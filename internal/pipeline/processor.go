package pipeline

import (
	"context"
	"db_migrate_server/internal/cast"
	"db_migrate_server/internal/expr"
	"db_migrate_server/internal/join"
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/sqlbuilder"
	"db_migrate_server/internal/util"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Processor struct {
	Plan  *mapping.Planner
	Pivot PivotStore
}

// NewProcessor wires planner and pivot repo into a Processor orchestrator.
func NewProcessor(plan *mapping.Planner, pivotRepo PivotStore) *Processor {
	return &Processor{Plan: plan, Pivot: pivotRepo}
}

// PivotStore membatasi operasi pivot yang boleh diakses processor.
type PivotStore interface {
	EnqueueNeedJoin(ctx context.Context, it pivot.NeedJoinItem) error
	AddJoinFragment(ctx context.Context, entity, joinKey, topic, sourceKey string, joinFields map[string]string, payload *kafka.DebeziumValue) error
	LookupKey(ctx context.Context, mapName, srcKey string) (string, bool, error)
	Enqueue(ctx context.Context, it pivot.ExecItem) error
	EnqueueSplit(ctx context.Context, queueID uuid.UUID, sqlText string, args []interface{}, returning []string, keymapPayload []byte) error
}

// KeyResolution captures resolved key value and optional keymap request info.
type KeyResolution struct {
	Value      string
	NeedKeymap bool // default false
	Request    *pivot.KeymapRequest
}

// Handle routes a single Kafka event to matching entities, deciding between
// multi-topic join path (store fragments or enqueue placeholders) and direct
// single-topic planning.
func (p *Processor) Handle(ctx context.Context, ev kafka.Event) error {

	// get matched entities with hash map
	entities := p.Plan.TopicToEntities[ev.Topic]
	util.Debug.Printf("processor: Handle topic=%s op=%s matchedEntities=%d", ev.Topic, ev.Op, len(entities))

	if len(entities) == 0 {
		util.Debug.Printf("no entity matched topic=%s", ev.Topic)
		return nil
	}

	for _, ent := range entities {

		util.Debug.Printf("processor: handling entity=%s op=%s", ent.Entity, ev.Op)

		// check fact condition
		if !CheckFactCondition(ev, &ent) {
			util.Debug.Printf("processor: skip entity=%s because fact_condition not met", ent.Entity)
			continue
		}

		// join wait if entity need more than 1 topic
		allTopics := mapping.ExpectedTopics(ent) // check if entity has multiple sources
		util.Debug.Printf("processor: expected topics=%v for entity=%s", allTopics, ent.Entity)

		if len(allTopics) > 1 { // if need mutiple sources / join

			aliasTopic := mapping.AliasToTopicHashMap(ent)

			// check if event is fact or dim then derive join key
			joinCtx := join.DeriveContext(ent, ev.Topic, ev.Value, aliasTopic)
			joinKey := joinCtx.JoinKey
			util.Debug.Printf("processor: join key=%s topic=%s fields=%v", joinKey, ev.Topic, joinCtx.Fields)

			// get value source key for keymap (if needed)
			sourceKey := p.deriveKeySource(ent, map[string]*kafka.DebeziumValue{ev.Topic: ev.Value})
			if sourceKey == "" {
				sourceKey = joinKey
			}
			util.Debug.Printf("processor: source key=%s topic=%s", sourceKey, ev.Topic)

			factTopic := topicName(ent.Sources[0])
			if ev.Topic == factTopic {
				// store need join row if fact event
				raw, _ := json.Marshal(ev.Value)
				fieldsRaw, _ := json.Marshal(joinCtx.Fields)
				item := pivot.NeedJoinItem{
					Entity:      ent.Entity,
					Op:          ev.Op,
					JoinKey:     joinKey,
					JoinTopic:   ev.Topic,
					JoinPayload: raw,
					JoinFields:  fieldsRaw,
				}
				if err := p.Pivot.EnqueueNeedJoin(ctx, item); err != nil {
					return err
				}
				util.Debug.Printf("processor: enqueued join placeholder entity=%s joinKey=%s", ent.Entity, joinKey)

				// store join fragment if dim event
			} else {
				if err := p.Pivot.AddJoinFragment(ctx, ent.Entity, joinKey, ev.Topic, sourceKey, joinCtx.Fields, ev.Value); err != nil {
					return err
				}
				util.Debug.Printf("processor: stored join fragment entity=%s topic=%s", ent.Entity, ev.Topic)
			}

		} else {
			util.Debug.Printf("processor: single topic path entity=%s topic=%s", ent.Entity, ev.Topic)

			if err := p.planAndEnqueue(ctx, ent, ev.Op, map[string]*kafka.DebeziumValue{ev.Topic: ev.Value}); err != nil {
				return err
			}
		}

	}

	return nil
}

// mergePayload merges payload from multiple topics (identity map here).
func (p *Processor) mergePayload(topicBytes map[string]*kafka.DebeziumValue) map[string]*kafka.DebeziumValue {
	return topicBytes
}

// planAndEnqueue resolves keys, chooses route, builds SQL, and writes to pivot queue.
func (p *Processor) planAndEnqueue(ctx context.Context, ent mapping.Entity, op string, payload map[string]*kafka.DebeziumValue) error {
	return p.planAndEnqueueWithQueue(ctx, ent, op, uuid.Nil, payload)
}

// planAndEnqueueWithQueue sama seperti planAndEnqueue, tapi meneruskan queueID jika ada (untuk joiner).
func (p *Processor) planAndEnqueueWithQueue(ctx context.Context, ent mapping.Entity, op string, queueID uuid.UUID, payload map[string]*kafka.DebeziumValue) error {

	util.Debug.Printf("processor: planAndEnqueue entity=%s op=%s payloadTopics=%d", ent.Entity, op, len(payload))

	// check foreign type
	keyRes, err := p.resolveKey(ctx, ent, payload)
	if err != nil {
		return err
	}
	util.Debug.Printf("processor: resolved key value=%s needKeymap=%t entity=%s", keyRes.Value, keyRes.NeedKeymap, ent.Entity)

	// lihat mode routing
	route := opRoute(op, ent)
	util.Debug.Printf("processor: route=%s entity=%s", route, ent.Entity)

	// 3) Bangun kolom target menyesuaikan mode penulisan
	cols, vals, sets, where, err := p.buildColumns(ent, keyRes.Value, keyRes.Value != "", payload, route)
	if err != nil {
		return err
	}
	util.Debug.Printf("processor: built columns entity=%s cols=%v", ent.Entity, cols)

	// 4) Tentukan writeMode & buat SQL
	var stmt sqlbuilder.Stmt

	switch route {
	case "insert":
		stmt = sqlbuilder.Insert(ent.TargetTable, cols, vals)
	case "update":
		if len(where) == 0 {
			return fmt.Errorf("missing match key for update route on entity %s", ent.Entity)
		}
		stmt = sqlbuilder.Update(ent.TargetTable, sets, where, vals)
	default:
		return fmt.Errorf("unsupported route mode %s for entity %s", route, ent.Entity)
	}

	// 5) Enqueue ke pivot._exec_queue
	var returning []string
	if keyRes.NeedKeymap && keyRes.Request != nil && keyRes.Request.TgtColumn != "" {
		stmt = stmt.WithReturning(keyRes.Request.TgtColumn)
		returning = []string{keyRes.Request.TgtColumn}
	}

	var keymapReq *pivot.KeymapRequest
	if keyRes.NeedKeymap {
		keymapReq = keyRes.Request
	}

	if err := p.Pivot.Enqueue(ctx, pivot.ExecItem{
		Entity:     ent.Entity,
		Op:         op,
		SQL:        stmt.SQL,
		Args:       stmt.Args,
		QueueID:    queueID,
		NeedKeymap: keyRes.NeedKeymap,
		Keymap:     keymapReq,
		Returning:  returning,
	}); err != nil {
		return err
	}

	// split_table: setiap target tambahan ditulis sebagai insert terpisah dengan queue_id yang sama
	for _, split := range ent.SplitTables {
		colsSplit, valsSplit, err := evalColumns(ent, keyRes.Value, keyRes.Value != "", payload, split.Columns, true)
		if err != nil {
			return err
		}
		stmtSplit := buildInsertStmt(split.TableName, colsSplit, valsSplit)
		var kmPayload []byte
		if keyRes.NeedKeymap && keyRes.Request != nil {
			if b, err := json.Marshal(keyRes.Request); err == nil {
				kmPayload = b
			}
		}
		if err := p.Pivot.EnqueueSplit(ctx, queueID, stmtSplit.SQL, stmtSplit.Args, nil, kmPayload); err != nil {
			return err
		}
	}

	return nil
}

// opRoute maps Debezium op to configured routing mode (insert/update).
func opRoute(op string, ent mapping.Entity) string {
	switch op {
	case "c":
		return normalizeRoute(ent.Routing.OnCreate.Mode, "insert")
	case "u":
		return normalizeRoute(ent.Routing.OnUpdate.Mode, "update")
	case "r":
		return normalizeRoute(ent.Routing.OnSnapshot.Mode, "insert")
	default:
		// fallback ke mode insert agar event tak dikenal tetap ditulis
		return normalizeRoute(ent.Routing.OnCreate.Mode, "insert")
	}
}

// normalizeRoute picks valid route or falls back when empty.
func normalizeRoute(mode, fallback string) string {
	switch mode {
	case "insert", "update":
		return mode
	case "":
		return fallback
	default:
		return mode
	}
}

// resolveKey implements key strategy: derive natural key or request/lookup keymap.
func (p *Processor) resolveKey(ctx context.Context, ent mapping.Entity, payload map[string]*kafka.DebeziumValue) (KeyResolution, error) {
	util.Debug.Printf("processor: resolveKey entity=%s strategy=%s", ent.Entity, ent.Key.Strategy)

	switch ent.Key.Strategy {

	case "natural":

		key := p.deriveKeySource(ent, payload)

		return KeyResolution{Value: key}, nil

	case "shared_key":
		if ent.Key.Resolver == nil || ent.Key.Resolver.Table == "" {
			return KeyResolution{}, fmt.Errorf("missing resolver for %s", ent.Entity)
		}
		mapName := ent.Key.Resolver.Table // atau pakai nama entitas (bebas)
		srcKey := p.deriveKeySource(ent, payload)
		if srcKey == "" {
			return KeyResolution{}, fmt.Errorf("missing source key for %s", ent.Entity)
		}

		// lookup
		if tgt, ok, err := p.Pivot.LookupKey(ctx, mapName, srcKey); err != nil {
			return KeyResolution{}, err
		} else if ok {
			return KeyResolution{Value: tgt}, nil
		}

		req := &pivot.KeymapRequest{
			MapName:   mapName,
			SrcKey:    srcKey,
			SrcColumn: sourceKeyColumn(ent),
			TgtColumn: targetKeyColumn(ent),
			SrcTable:  firstSourceTable(ent),
			TgtTable:  ent.TargetTable,
		}
		return KeyResolution{NeedKeymap: true, Request: req}, nil
	default:
		return KeyResolution{}, fmt.Errorf("unknown key strategy: %s (supported: natural, shared_key)", ent.Key.Strategy)
	}
}

// anyPayload returns an arbitrary payload from the map (when only one is expected).
func anyPayload(m map[string]*kafka.DebeziumValue) *kafka.DebeziumValue {
	for _, v := range m {
		util.Debug.Printf("processor: anyPayload returning payload topic")

		return v
	}
	util.Debug.Printf("processor: anyPayload no payload available")
	return nil
}

// stringFromRow tries to stringify a column value from Debezium row map.
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

// sourceKeyColumn resolves column name used as src_key for keymap requests.
func sourceKeyColumn(ent mapping.Entity) string {
	if ent.Key.Resolver != nil && ent.Key.Resolver.SourceKeyCol != "" {
		return ent.Key.Resolver.SourceKeyCol
	}
	if strings.TrimSpace(ent.Key.Source) != "" {
		parts := strings.Split(ent.Key.Source, ".")
		return parts[len(parts)-1]
	}
	return ""
}

// targetKeyColumn finds target column bound to $key placeholder.
func targetKeyColumn(ent mapping.Entity) string {
	for colName, spec := range ent.Columns {
		if spec.From == "$key" {
			return colName
		}
	}
	return ""
}

// firstSourceTable returns the first source table name for metadata/debug.
func firstSourceTable(ent mapping.Entity) string {
	if len(ent.Sources) == 0 {
		return ""
	}
	return ent.Sources[0].From
}

// topicName derives topic from explicit field or table name.
func topicName(s mapping.EntitySource) string {
	if trimmed := strings.TrimSpace(s.Topic); trimmed != "" {
		return trimmed
	}
	name := strings.TrimSpace(s.From)
	if name == "" {
		return ""
	}
	parts := strings.Split(name, ".")
	base := parts[len(parts)-1]
	base = strings.ToLower(strings.TrimSpace(base))
	if base == "" {
		return ""
	}
	base = strings.ReplaceAll(base, " ", "_")
	return fmt.Sprintf("db_events_%s", base)
}

// deriveKeySource mengambil source key sesuai key.source dan payload gabungan (topic->payload).
// Ini dipakai untuk keymap (shared_key) agar tidak bergantung pada join key.
func (p *Processor) deriveKeySource(ent mapping.Entity, payload map[string]*kafka.DebeziumValue) string {
	aliasTopic := mapping.AliasToTopicHashMap(ent)

	src := strings.TrimSpace(ent.Key.Source)
	if src == "" {
		return ""
	}
	alias := strings.Split(src, ".")[0]

	topic := aliasTopic[alias]

	return keyFromSource(src, "", payload[topic], aliasTopic)
}

// HandleJoinReady dipakai checker saat join fragment sudah lengkap.
func (p *Processor) HandleJoinReady(ctx context.Context, ent mapping.Entity, op string, payload map[string]*kafka.DebeziumValue) error {
	return p.planAndEnqueue(ctx, ent, op, payload)
}

// HandleJoinReadyWithQueue memungkinkan joiner meneruskan queueID agar konsisten dengan _need_join.
func (p *Processor) HandleJoinReadyWithQueue(ctx context.Context, ent mapping.Entity, op string, queueID uuid.UUID, payload map[string]*kafka.DebeziumValue) error {
	return p.planAndEnqueueWithQueue(ctx, ent, op, queueID, payload)
}

// keyFromSource mengambil nilai kolom dari key.source dengan memperhatikan alias->topic.
func keyFromSource(src string, eventTopic string, value *kafka.DebeziumValue, aliasTopic map[string]string) string {
	src = strings.TrimSpace(src)
	if value == nil || src == "" {
		return ""
	}
	parts := strings.Split(src, ".")
	alias := parts[0]
	col := parts[len(parts)-1]
	topic := aliasTopic[alias]
	if topic != "" && eventTopic != "" && topic != eventTopic {
		return ""
	}
	return stringFromRow(value.Payload.After, col)
}

// buildColumns assembles target columns/values and where clause for insert/update.
func (p *Processor) buildColumns(ent mapping.Entity, key string, hasKey bool, payload map[string]*kafka.DebeziumValue, route string) (
	cols []string, vals []interface{}, sets []string, where []string, err error) {

	// Disederhanakan: hanya handle "$key" dan "from: <col>" dari after.*; dukungan expr terbatas ke fungsi built-in.
	// Kamu bisa perluasan: cast, resolver lookup, flatten/aggregate (untuk contoh ringkas ini belum penuh).

	for colName, spec := range ent.Columns {
		var (
			value   interface{}
			handled bool
		)
		if spec.Expr != "" {
			if eval, ok, evalErr := expr.Evaluate(spec.Expr); evalErr != nil {
				return cols, vals, sets, nil, fmt.Errorf("entity %s column %s: %w", ent.Entity, colName, evalErr)
			} else if ok {
				value = eval
				handled = true
			}
		}

		switch {
		case handled:
			// value already populated by expression
		case spec.From == "$key":
			if !hasKey {
				util.Debug.Printf("processor: skipping column=%s entity=%s because key unavailable", colName, ent.Entity)
				continue
			}
			value = key
		case spec.From != "":
			value = extractFieldFromPayload(payload, spec.From)
		default:
			value = spec.Default
		}

		value, err = cast.Value(spec.Cast, value)
		if err != nil {
			return cols, vals, sets, nil, fmt.Errorf("entity %s column %s: %w", ent.Entity, colName, err)
		}
		if spec.Cast != "" {
			util.Debug.Printf("processor: cast column=%s entity=%s cast=%s -> %v", colName, ent.Entity, spec.Cast, value)
		}

		cols = append(cols, colName)
		vals = append(vals, value)
	}

	// sets untuk UPDATE (col=$i), where dari matchKey atau primary (key)
	for i, c := range cols {
		sets = append(sets, fmt.Sprintf("%s=$%d", c, i+1))
	}

	if route == "update" {
		match := ent.Routing.OnUpdate.MatchKey
		if len(match) == 0 {
			if tk := targetKeyColumn(ent); tk != "" {
				match = []string{tk}
			}
		}
		if len(match) == 0 {
			return cols, vals, sets, nil, fmt.Errorf("missing match key for update entity=%s", ent.Entity)
		}
		indexes := map[string]int{}
		for idx, c := range cols {
			indexes[c] = idx + 1 // placeholder index
		}
		for _, k := range match {
			if pos, ok := indexes[k]; ok {
				where = append(where, fmt.Sprintf("%s=$%d", k, pos))
			} else {
				return cols, vals, sets, nil, fmt.Errorf("match key %s not present in columns for entity %s", k, ent.Entity)
			}
		}
	}

	return
}

// extractFieldFromPayload picks a column value from merged payload using the last path segment.
func extractFieldFromPayload(payload map[string]*kafka.DebeziumValue, path string) any {
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

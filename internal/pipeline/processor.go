package pipeline

import (
	"context"
	"db_migrate_server/internal/cast"
	"db_migrate_server/internal/expr"
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/sqlbuilder"
	"db_migrate_server/internal/util"
	"encoding/json"
	"fmt"
	"strings"
)

type Processor struct {
	Plan  *mapping.Planner
	Pivot *pivot.Repo
}

func NewProcessor(plan *mapping.Planner, pivotRepo *pivot.Repo) *Processor {
	return &Processor{Plan: plan, Pivot: pivotRepo}
}

type KeyResolution struct {
	Value      string
	NeedKeymap bool // default false
	Request    *pivot.KeymapRequest
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

			joinKey := p.deriveJoinKey(ent, ev.Topic, ev.Value)
			util.Debug.Printf("processor: join key=%s topic=%s", joinKey, ev.Topic)

			sourceKey := p.deriveKeySource(ent, map[string]*kafka.DebeziumValue{ev.Topic: ev.Value})
			if sourceKey == "" {
				sourceKey = joinKey
			}
			util.Debug.Printf("processor: source key=%s topic=%s", sourceKey, ev.Topic)

			factTopic := topicName(ent.Sources[0])
			if ev.Topic == factTopic {
				// enqueue placeholder needing join
				raw, _ := json.Marshal(ev.Value)
				item := pivot.ExecItem{
					Entity:        ent.Entity,
					Op:            op,
					NeedJoin:      true,
					JoinKey:       joinKey,
					JoinTopic:     ev.Topic,
					JoinSourceKey: sourceKey,
					JoinPayload:   raw,
					NeedKeymap:    false,
					SQL:           "pending_join",
					Args:          nil,
					Returning:     nil,
				}
				if err := p.Pivot.Enqueue(ctx, item); err != nil {
					return err
				}
				util.Debug.Printf("processor: enqueued join placeholder entity=%s joinKey=%s", ent.Entity, joinKey)
			} else {
				if err := p.Pivot.AddJoinFragment(ctx, ent.Entity, joinKey, ev.Topic, sourceKey, ev.Value); err != nil {
					return err
				}
				util.Debug.Printf("processor: stored join fragment entity=%s topic=%s", ent.Entity, ev.Topic)
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

func (p *Processor) deriveJoinKey(ent mapping.Entity, eventTopic string, value *kafka.DebeziumValue) string {

	aliasTopic := aliasTopics(ent)

	// Ambil dari ekspresi join "alias.col = alias2.col" yang sesuai topic event
	if key := keyFromJoinHint(ent, eventTopic, value, aliasTopic); key != "" {
		return key
	}

	// Fallback: coba key.source jika join expr tidak cocok
	if len(ent.Key.Source) > 0 {
		if key := keyFromSource(ent.Key.Source, eventTopic, value, aliasTopic); key != "" {
			util.Debug.Printf("processor: derived join key=%s entity=%s via key.source", key, ent.Entity)
			return key
		}
	}

	util.Debug.Printf("processor: deriveJoinKey missing value entity=%s", ent.Entity)
	return ""
}

// keyFromJoinHint mencoba membaca ent.Key.JoinKey terlebih dulu; bila tidak ada, pakai join expression.
func keyFromJoinHint(ent mapping.Entity, eventTopic string, value *kafka.DebeziumValue, aliasTopic map[string]string) string {
	if value == nil {
		return ""
	}

	if ent.Key.JoinKey != "" {
		if key := keyFromSource([]string{ent.Key.JoinKey}, eventTopic, value, aliasTopic); key != "" {
			util.Debug.Printf("processor: derived join key=%s entity=%s via join_key", key, ent.Entity)
			return key
		}
	}

	for _, src := range ent.Sources {
		if src.Join == nil || strings.TrimSpace(src.Join.On) == "" {
			continue
		}
		clauses := strings.Split(strings.ReplaceAll(src.Join.On, "AND", "and"), "and")
		for _, c := range clauses {
			partsEq := strings.Split(c, "=")
			if len(partsEq) != 2 {
				continue
			}
			left := strings.TrimSpace(partsEq[0])
			right := strings.TrimSpace(partsEq[1])
			for _, side := range []string{left, right} {
				ps := strings.Split(side, ".")
				if len(ps) < 2 {
					continue
				}
				alias := strings.TrimSpace(ps[0])
				column := strings.TrimSpace(ps[len(ps)-1])
				topic := aliasTopic[alias]
				if topic == "" || (eventTopic != "" && topic != eventTopic) {
					continue
				}
				key := stringFromRow(value.Payload.After, column)
				if key == "" {
					key = stringFromRow(value.Payload.Before, column)
				}
				if key != "" {
					util.Debug.Printf("processor: derived join key=%s entity=%s via join expr", key, ent.Entity)
					return key
				}
			}
		}
	}
	return ""
}

// mergePayload: gabungkan payload dari beberapa topic (sederhana: map[topic]rawJSON)
func (p *Processor) mergePayload(topicBytes map[string]*kafka.DebeziumValue) map[string]*kafka.DebeziumValue {
	return topicBytes
}

func (p *Processor) planAndEnqueue(ctx context.Context, ent mapping.Entity, op string, payload map[string]*kafka.DebeziumValue) error {

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

	return p.Pivot.Enqueue(ctx, pivot.ExecItem{
		Entity: ent.Entity,
		Op:     op,
		SQL:    stmt.SQL,
		Args:   stmt.Args,

		NeedKeymap: keyRes.NeedKeymap,
		Keymap:     keymapReq,
		Returning:  returning,
	})
}

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

// mendefinisikan mapping foreign berdasarkan 'strategy'
func (p *Processor) resolveKey(ctx context.Context, ent mapping.Entity, payload map[string]*kafka.DebeziumValue) (KeyResolution, error) {
	util.Debug.Printf("processor: resolveKey entity=%s strategy=%s", ent.Entity, ent.Key.Strategy)

	switch ent.Key.Strategy {

	case "natural":

		key := p.deriveKeySource(ent, payload)

		return KeyResolution{Value: key}, nil

	case "shared_key", "surrogate":
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
		return KeyResolution{}, fmt.Errorf("unknown key strategy")
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

func sourceKeyColumn(ent mapping.Entity) string {
	if ent.Key.Resolver != nil && ent.Key.Resolver.SourceKeyCol != "" {
		return ent.Key.Resolver.SourceKeyCol
	}
	if len(ent.Key.Source) > 0 {
		parts := strings.Split(ent.Key.Source[0], ".")
		return parts[len(parts)-1]
	}
	return ""
}

func targetKeyColumn(ent mapping.Entity) string {
	for colName, spec := range ent.Columns {
		if spec.From == "$key" {
			return colName
		}
	}
	return ""
}

func firstSourceTable(ent mapping.Entity) string {
	if len(ent.Sources) == 0 {
		return ""
	}
	return ent.Sources[0].From
}

func aliasTopics(ent mapping.Entity) map[string]string {
	out := map[string]string{}
	for _, s := range ent.Sources {
		topic := topicName(s)
		if s.Alias != "" && topic != "" {
			out[s.Alias] = topic
		}
	}
	return out
}

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
// Ini dipakai untuk keymap (surrogate/shared) agar tidak bergantung pada join key.
func (p *Processor) deriveKeySource(ent mapping.Entity, payload map[string]*kafka.DebeziumValue) string {
	aliasTopic := aliasTopics(ent)
	for _, src := range ent.Key.Source {
		if key := keyFromSource([]string{src}, "", payload[aliasTopic[strings.Split(src, ".")[0]]], aliasTopic); key != "" {
			return key
		}
	}
	return ""
}

// HandleJoinReady dipakai checker saat join fragment sudah lengkap.
func (p *Processor) HandleJoinReady(ctx context.Context, ent mapping.Entity, op string, payload map[string]*kafka.DebeziumValue) error {
	return p.planAndEnqueue(ctx, ent, op, payload)
}

// keyFromSource mengambil nilai kolom dari key.source dengan memperhatikan alias->topic
func keyFromSource(keySources []string, eventTopic string, value *kafka.DebeziumValue, aliasTopic map[string]string) string {
	if value == nil || len(keySources) == 0 {
		return ""
	}
	for _, src := range keySources {
		parts := strings.Split(src, ".")
		if len(parts) == 0 {
			continue
		}
		alias := parts[0]
		col := parts[len(parts)-1]
		topic := aliasTopic[alias]
		if topic != "" && eventTopic != "" && topic != eventTopic {
			continue
		}
		if key := stringFromRow(value.Payload.After, col); key != "" {
			return key
		}
		if key := stringFromRow(value.Payload.Before, col); key != "" {
			return key
		}
	}
	return ""
}

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
			value = firstAfterField(payload, spec.From)
		default:
			value = spec.Default
		}

		value, err = cast.Value(spec.Cast, value)
		if err != nil {
			return cols, vals, sets, nil, fmt.Errorf("entity %s column %s: %w", ent.Entity, colName, err)
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

// ambil nama kolom asli tanpa path
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

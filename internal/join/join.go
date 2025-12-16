package join

import (
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/util"
	"fmt"
	"sort"
	"strings"
)

// Context merangkum hasil derivasi join dari satu event.
type Context struct {
	JoinKey string
	Fields  map[string]string
	IsFact  bool
}

type spec struct {
	FactAlias string
	FactCol   string
	DimAlias  string
	DimCol    string
}

// DeriveContext picks join field values for the current topic, producing the
// composite join key, the field map used, dan flag apakah event ini fact.
func DeriveContext(ent mapping.Entity, eventTopic string, value *kafka.DebeziumValue, aliasTopic map[string]string) Context {

	specs := collectSpecs(ent) // create join fk mappinng from this entity
	// example if join 3 tables (need 2 dim table): [
	// {factAlias:"p" factCol:"pub_cntry" dimAlias:"c" dimCol:"cntry_id"},
	// {factAlias:"p" factCol:"pub_conf" dimAlias:"f" dimCol:"conf_id"}
	// ...]

	fields := map[string]string{}
	isFact := false

	// bruto force check each join spec (nested loop)
	for _, spec := range specs {
		currentIsFact := aliasTopic[spec.FactAlias] == eventTopic
		currentIsDim := aliasTopic[spec.DimAlias] == eventTopic
		if !currentIsFact && !currentIsDim {
			continue
		}

		if currentIsFact {
			isFact = true
		}

		col := spec.DimCol // if dim, pakai col dim

		if currentIsFact { // if fact, pakai col fact
			col = spec.FactCol
		}

		val := valueFromPayload(value, col) // get fk value from payload
		if val == "" {
			continue
		}
		key := fmt.Sprintf("%s.%s::%s.%s", spec.FactAlias, spec.FactCol, spec.DimAlias, spec.DimCol)
		fields[key] = val
		// example key: "p.pub_cntry::c.cntry_id" value:"DE"
	}

	joinKey := joinKeyFromFields(fields)
	// example value of joinKey if dim event come --> "p.pub_id::f.conf_id=14"
	// example value of joinKey if fact event come --> "p.pub_cntry::c.cntry_id=ES|p.pub_id::f.conf_id=2"

	if joinKey == "" {
		util.Debug.Printf("join: derive missing value entity=%s topic=%s", ent.Entity, eventTopic)
	}
	return Context{JoinKey: joinKey, Fields: fields, IsFact: isFact}
}

// joinKeyFromFields builds a deterministic composite join key from field-value pairs.
func joinKeyFromFields(fields map[string]string) string {
	if len(fields) == 0 {
		return ""
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, fields[k]))
	}
	return strings.Join(parts, "|")
}

// collectSpecs rewrites mapping join config into internal spec list.
func collectSpecs(ent mapping.Entity) []spec {
	var specs []spec

	fact := ""
	if len(ent.Sources) > 0 {
		fact = ent.Sources[0].Alias
	}

	for _, s := range ent.Sources {
		if s.Join == nil {
			continue
		}

		specs = append(specs, spec{
			FactAlias: fact,
			FactCol:   s.Join.FactColumn,
			DimAlias:  s.Alias,
			DimCol:    s.Join.DimColumn,
		})
	}

	return specs
}

// valueFromPayload extracts a column value from Debezium after/before maps.
func valueFromPayload(row *kafka.DebeziumValue, column string) string {
	if row == nil || column == "" {
		return ""
	}
	if v := stringFromRow(row.Payload.After, column); v != "" {
		return v
	}
	return stringFromRow(row.Payload.Before, column)
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

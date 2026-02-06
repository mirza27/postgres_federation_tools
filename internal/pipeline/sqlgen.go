package pipeline

import (
	"db_migrate_server/internal/cast"
	"db_migrate_server/internal/expr"
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/sqlbuilder"
	"db_migrate_server/internal/util"
	"fmt"
	"strings"
)

const splitKeyPlaceholder = "__KEYMAP_PLACEHOLDER__"

// buildInsertStmt menyusun INSERT sederhana berdasarkan kolom yang sudah dievaluasi.
func buildInsertStmt(table string, cols []string, vals []interface{}) sqlbuilder.Stmt {
	return sqlbuilder.Insert(table, cols, vals)
}

// evalColumns mengevaluasi spesifikasi kolom (from/expr/default/$key) untuk payload tertentu.
// Dipakai baik untuk target utama maupun split_table.
func evalColumns(ent mapping.Entity, key string, hasKey bool, payload map[string]*kafka.DebeziumValue, specs map[string]mapping.Column, usePlaceholder bool) (cols []string, vals []interface{}, err error) {
	for colName, spec := range specs {
		var (
			value   interface{}
			handled bool
		)
		if spec.Expr != "" {
			if eval, ok, evalErr := expr.Evaluate(spec.Expr); evalErr != nil {
				return cols, vals, fmt.Errorf("entity %s column %s: %w", ent.Entity, colName, evalErr)
			} else if ok {
				value = eval
				handled = true
			}
		}
		switch {
		case handled:
		case spec.From == "$key":
			if !hasKey {
				if usePlaceholder {
					value = splitKeyPlaceholder
				} else {
					util.Debug.Printf("sqlgen: skip column=%s entity=%s because key unavailable", colName, ent.Entity)
					continue
				}
			} else {
				value = key
			}
		case spec.From != "":
			value = firstAfterField(payload, spec.From)
		default:
			value = spec.Default
		}
		value, err = cast.Value(spec.Cast, value)
		if err != nil {
			return cols, vals, fmt.Errorf("entity %s column %s: %w", ent.Entity, colName, err)
		}
		if spec.Cast != "" {
			util.Debug.Printf("sqlgen: cast column=%s entity=%s split=%t cast=%s -> %v", colName, ent.Entity, usePlaceholder, spec.Cast, value)
		}
		cols = append(cols, colName)
		vals = append(vals, value)
	}
	return
}

// firstAfterField picks a column value from merged payload using the last path segment.
func firstAfterField(payload map[string]*kafka.DebeziumValue, path string) any {
	util.Debug.Printf("sqlgen: firstAfterField path=%s", path)
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

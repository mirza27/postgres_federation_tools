package pipeline

import (
	"db_migrate_server/internal/kafka"
	"db_migrate_server/internal/mapping"
	"fmt"
	"strings"
)

// return true jika sebuah event fact memenuhi condition atau jika sebuah bagian event dimensi
func CheckFactCondition(ev kafka.Event, ent *mapping.Entity) bool {

	if ent.FactCondition == nil {
		return true
	}

	// check apakah event adalah bagian fact
	factTopic := mapping.SourceTopic(ent.Sources[0])
	if ev.Topic != factTopic {

		// still process non-fact event
		return true
	}

	// cek apakah kondisi terpenuhi
	return factConditionMatch(ent.FactCondition, ev.Value)
}

func factConditionMatch(fc *mapping.FactCondition, value *kafka.DebeziumValue) bool {
	if fc == nil || value == nil {
		return true
	}
	colVal := stringFromRow(value.Payload.After, fc.Column)
	if colVal == "" {
		colVal = stringFromRow(value.Payload.Before, fc.Column)
	}
	target := fmt.Sprint(fc.Value)
	switch strings.ToLower(fc.Op) {
	case "equal":
		return colVal == target
	case "notequals":
		return colVal != target
	default:
		// op tak dikenal dianggap lolos untuk backward compatibility
		return true
	}
}

package kafka

import (
	"encoding/json"
	"fmt"
)

// Event represents a single Debezium event consumed from Kafka.
// The value payload is decoded into structured fields so downstream
// processors can access data without reparsing raw bytes.
// pemetaan yang digunakan sistem
type Event struct {
	Topic  string
	Key    []byte
	Op     string
	Source *DebeziumSource
	Value  *DebeziumValue
}

// DebeziumValue mirrors the useful portions of the Debezium envelope.
// Only a subset of fields is captured for the processor's needs.
type DebeziumValue struct {
	Payload DebeziumPayload
	Op      string // add op for easy access
}

// DebeziumPayload contains the before/after rows and the source metadata.
type DebeziumPayload struct {
	Before map[string]any  `json:"before"`
	After  map[string]any  `json:"after"`
	Source *DebeziumSource `json:"source"`
	Op     string          `json:"op"`
}

// DebeziumSource captures metadata about the event origin.
type DebeziumSource struct {
	Version   string  `json:"version"`
	Connector string  `json:"connector"`
	Name      string  `json:"name"`
	TsMs      *int64  `json:"ts_ms"`
	Snapshot  string  `json:"snapshot"`
	Db        string  `json:"db"`
	Sequence  *string `json:"sequence"`
	Schema    string  `json:"schema"`
	Table     string  `json:"table"`
	TxID      *int64  `json:"txId"`
	Lsn       *int64  `json:"lsn"`
	Xmin      *int64  `json:"xmin"`
}

// DecodeEventValue decodes the raw Debezium envelope into structured fields.
func DecodeEventValue(raw []byte) (*DebeziumValue, error) {
	// get value from debezium with key "payload"
	var tmp struct {
		Payload *DebeziumPayload `json:"payload"`
	}

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return nil, err
	}

	if tmp.Payload == nil {
		return &DebeziumValue{}, nil
	}

	payload := tmp.Payload
	debeziumValue := &DebeziumValue{
		Payload: *payload,
		Op:      payload.Op, // copy op for easy access
	}

	return debeziumValue, nil
}

// String renders the event in JSON form for debugging.
func (e Event) String() string {
	if e.Value == nil && e.Source == nil {
		return fmt.Sprintf("{topic:%q op:%q key:%x value:<nil> source:<nil>}", e.Topic, e.Op, e.Key)
	}
	type dbg struct {
		Topic  string          `json:"topic"`
		Key    []byte          `json:"key,omitempty"`
		Op     string          `json:"op,omitempty"`
		Source *DebeziumSource `json:"source,omitempty"`
		Value  *DebeziumValue  `json:"value,omitempty"`
	}
	d := dbg{
		Topic:  e.Topic,
		Key:    e.Key,
		Op:     e.Op,
		Source: e.Source,
		Value:  e.Value,
	}
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("{topic:%q op:%q key:%x value:<error:%v>}", e.Topic, e.Op, e.Key, err)
	}
	return string(b)
}

// String renders the DebeziumValue content as JSON for easier inspection.
func (v *DebeziumValue) String() string {
	if v == nil {
		return "<nil>"
	}
	type payload struct {
		Before map[string]any  `json:"before,omitempty"`
		After  map[string]any  `json:"after,omitempty"`
		Source *DebeziumSource `json:"source,omitempty"`
	}
	type dbg struct {
		Op      string  `json:"op,omitempty"`
		TsMs    *int64  `json:"ts_ms,omitempty"`
		Payload payload `json:"payload"`
	}
	d := dbg{
		Op: v.Op,
		Payload: payload{
			Before: v.Payload.Before,
			After:  v.Payload.After,
			Source: v.Payload.Source,
		},
	}
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("<error:%v>", err)
	}
	return string(b)
}

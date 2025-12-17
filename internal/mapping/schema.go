package mapping

type Root struct {
	Version  string   `json:"version"`
	Entities []Entity `json:"entities"`
}

type Entity struct {
	Entity      string         `json:"entity"`
	Sources     []EntitySource `json:"sources"`
	TargetTable string         `json:"target_table"`

	Key           Key               `json:"key"`
	Columns       map[string]Column `json:"columns"`
	FactCondition *FactCondition    `json:"fact_condition,omitempty"`
	SplitTables   []SplitTable      `json:"split_table,omitempty"`

	Routing Routing `json:"routing"`
}

type FactCondition struct {
	Column string      `json:"column"`
	Op     string      `json:"op"`    // equal | notEquals
	Value  interface{} `json:"value"` // dibandingkan sebagai string
}

// SplitTable menggambarkan target tambahan yang ditulis bersamaan dengan target utama.
type SplitTable struct {
	TableName string            `json:"table_name"`
	Columns   map[string]Column `json:"columns"`
}

type EntitySource struct {
	Alias string `json:"alias"`
	From  string `json:"from"`
	Topic string `json:"topic"`
	Join  *Join  `json:"join,omitempty"`
}

type Join struct {
	FactColumn string `json:"fact_column,omitempty"`
	DimColumn  string `json:"dim_column,omitempty"`
}

type Key struct {
	Strategy string    `json:"strategy"` // natural | shared_key
	Source   string    `json:"source,omitempty"`
	JoinKey  string    `json:"join_key,omitempty"`
	Resolver *Resolver `json:"resolver,omitempty"`
}

type Resolver struct {
	Type         string `json:"type"` // mapping_table | mapping_table_lookup
	Table        string `json:"table"`
	SourceKeyCol string `json:"source_key_col"`
	TargetKeyCol string `json:"target_key_col"`
	ValueFrom    string `json:"value_from,omitempty"`
	Gen          *Gen   `json:"gen,omitempty"`
}

type Gen struct {
	Type string `json:"type"` // uuid_v7 dll.
}

type Column struct {
	From     string      `json:"from,omitempty"` // "u.email" / "$key"
	Expr     string      `json:"expr,omitempty"` // "concat(...)" / "now()"
	Cast     string      `json:"cast,omitempty"`
	Default  interface{} `json:"default,omitempty"`
	Resolver *Resolver   `json:"resolver,omitempty"` // lookup FK via keymap lain
}

type Routing struct {
	OnCreate   Route `json:"on_create"`
	OnUpdate   Route `json:"on_update"`
	OnSnapshot Route `json:"on_snapshot"`
}

type Route struct {
	Mode     string   `json:"mode"` // insert|update
	MatchKey []string `json:"matchKey,omitempty"`
}

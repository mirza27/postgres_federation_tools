package mapping

type Root struct {
	Version  string   `json:"version"`
	Sources  []Source `json:"sources"`
	Target   Target   `json:"target"`
	Topics   *Topics  `json:"topics,omitempty"`
	Engine   Engine   `json:"engine"`
	Entities []Entity `json:"entities"`
}

type Source struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	DefaultSchema string `json:"default_schema"`
}

type Target struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Schema string `json:"schema"`
}

type Topics struct {
	Enabled bool              `json:"enabled"`
	Map     map[string]string `json:"map"`
}

type Engine struct {
	Language   string `json:"language"`
	SQLDialect string `json:"sqlDialect"`
	Batch      struct {
		MaxRows       int `json:"maxRows"`
		MaxIntervalMs int `json:"maxIntervalMs"`
	} `json:"batch"`
	Logging struct {
		Enabled bool   `json:"enabled"`
		Table   string `json:"table"`
		Level   string `json:"level"`
	} `json:"logging"`
	PivotDb struct {
		Type string `json:"type"`
		DSN  string `json:"dsn"`
	} `json:"pivotDb"`
	Expr struct {
		Dialect string `json:"dialect"`
	} `json:"expr"`
}

type Entity struct {
	Entity      string         `json:"entity"`
	Kind        string         `json:"kind"`
	Sources     []EntitySource `json:"sources"`
	TargetTable string         `json:"target_table"`

	Key            Key               `json:"key"`
	Columns        map[string]Column `json:"columns"`
	FactCondition  *FactCondition    `json:"fact_condition,omitempty"`

	Flatten    *Flatten    `json:"flatten,omitempty"`
	Aggregates *Aggregates `json:"aggregates,omitempty"`

	Routing Routing `json:"routing"`
}

type FactCondition struct {
	Column string      `json:"column"`
	Op     string      `json:"op"`    // equal | notEquals
	Value  interface{} `json:"value"` // dibandingkan sebagai string
}

type EntitySource struct {
	Alias string `json:"alias"`
	From  string `json:"from"`
	Topic string `json:"topic"`
	Join  *Join  `json:"join,omitempty"`
}

type Join struct {
	MatchWith  string `json:"match_with,omitempty"`
	FactColumn string `json:"fact_column,omitempty"`
	DimColumn  string `json:"dim_column,omitempty"`
}

type Key struct {
	Strategy string    `json:"strategy"` // natural | shared_key | surrogate
	Source   []string  `json:"source,omitempty"`
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

type Flatten struct {
	FromArray string `json:"from_array"`
	As        string `json:"as"`
	RowKey    Column `json:"row_key"`
}

type Aggregates struct {
	GroupBy []string          `json:"group_by"`
	Metrics map[string]Column `json:"metrics"`
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

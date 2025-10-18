package models

type WalChange struct {
	Change []ChangeItem `json:"change"`
}

type ChangeItem struct {
	Kind         string        `json:"kind"`
	Schema       string        `json:"schema"`
	Table        string        `json:"table"`
	ColumnNames  []string      `json:"columnnames"`
	ColumnTypes  []string      `json:"columntypes"`
	ColumnValues []interface{} `json:"columnvalues"`
	OldKeys      *OldKeys      `json:"oldkeys,omitempty"`
}

type OldKeys struct {
	KeyNames  []string      `json:"keynames"`
	KeyTypes  []string      `json:"keytypes"`
	KeyValues []interface{} `json:"keyvalues"`
}

type Wal2JsonMessage struct {
	Change []Wal2JsonChange `json:"change"`
}

type Wal2JsonChange struct {
	Kind         string                 `json:"kind"`
	Schema       string                 `json:"schema"`
	Table        string                 `json:"table"`
	ColumnNames  []string               `json:"columnnames"`
	ColumnTypes  []string               `json:"columntypes"`
	ColumnValues []interface{}          `json:"columnvalues"`
	OldKeys      *Wal2JsonOldKeys       `json:"oldkeys,omitempty"`
}

type Wal2JsonOldKeys struct {
	KeyNames  []string      `json:"keynames"`
	KeyTypes  []string      `json:"keytypes"`
	KeyValues []interface{} `json:"keyvalues"`
}

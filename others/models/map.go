package models

type MappingConfig struct {
	Version            string                            `json:"version"`
	Defaults           map[string]map[string]interface{} `json:"defaults"`
	Transformations    map[string]string                 `json:"transformations"`
	Mappings           []TableMapping                    `json:"mappings"`
	ConflictResolution ConflictResolution                `json:"conflict_resolution"`
}

type TableMapping struct {
	SourceTable     string        `json:"source_table"`
	SourceCondition string        `json:"source_condition,omitempty"`
	TargetTables    []TargetTable `json:"target_tables"`
}

type TargetTable struct {
	Database          string            `json:"database"`
	Table             string            `json:"table"`
	OperationHandling OperationHandling `json:"operation_handling"`
	ColumnMappings    []ColumnMapping   `json:"column_mappings"`
	KeyMapping        KeyMapping        `json:"key_mapping,omitempty"`
}

type ColumnMapping struct {
	SourceColumn     string      `json:"source_column,omitempty"`
	SourceExpression string      `json:"source_expression,omitempty"`
	TargetColumn     string      `json:"target_column"`
	Transformation   string      `json:"transformation,omitempty"`
	DefaultValue     interface{} `json:"default_value,omitempty"`
}

type KeyMapping struct {
	PrimaryKey SourceTarget `json:"primary_key,omitempty"`
}

type SourceTarget struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type OperationHandling struct {
	Insert bool   `json:"insert"`
	Update bool   `json:"update"`
	Delete string `json:"delete"` // "hard" or "soft_delete"
}

type ConflictResolution struct {
	Strategy        string `json:"strategy"`
	TimestampColumn string `json:"timestamp_column"`
}

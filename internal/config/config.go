package config

type DatabaseCredential struct {
	Type   string
	Host   string
	Port   int
	User   string
	Pass   string
	DbName string
}

type Config struct {
	TargetDatabaseType string `mapstructure:"TARGET_TYPE"`
	TargetDatabaseName string `mapstructure:"TARGET_DATABASE"`
	TargetDatabaseHost string `mapstructure:"TARGET_HOST"`
	TargetDatabasePort int    `mapstructure:"TARGET_PORT"`
	TargetDatabaseUser string `mapstructure:"TARGET_USER"`
	TargetDatabasePass string `mapstructure:"TARGET_PASSWORD"`

	SourceDatabaseType string `mapstructure:"SOURCE_TYPE"`
	SourceDatabaseName string `mapstructure:"SOURCE_DATABASE"`
	SourceDatabaseHost string `mapstructure:"SOURCE_HOST"`
	SourceDatabasePort int    `mapstructure:"SOURCE_PORT"`
	SourceDatabaseUser string `mapstructure:"SOURCE_USER"`
	SourceDatabasePass string `mapstructure:"SOURCE_PASSWORD"`

	TargetDSN       string `mapstructure:"TARGET_DSN"`
	PivotDSN        string `mapstructure:"PIVOT_DSN"`
	PivotSchemaPath string `mapstructure:"PIVOT_SCHEMA_PATH"`
	MappingPath     string `mapstructure:"MAPPING_PATH"`

	DebeziumHost          string `mapstructure:"DEBEZIUM_HOST"`
	DebeziumPort          int    `mapstructure:"DEBEZIUM_PORT"`
	DebeziumConnectorName string `mapstructure:"DEBEZIUM_CONNECTOR_NAME"`

	BatchMaxRows        int    `mapstructure:"BATCH_MAX_ROWS"`
	BatchMaxInterval    int    `mapstructure:"BATCH_MAX_INTERVAL"`
	MaximumJoinAttempts int    `mapstructure:"MAX_JOIN_ATTEMPTS"`
	KafkaBrokers        string `mapstructure:"KAFKA_BROKERS"`
	KafkaGroupID        string `mapstructure:"KAFKA_GROUP_ID"`
	KafkaCounsumerReset bool   `mapstructure:"KAFKA_CONSUMER_RESET"`

	RunApiPort int `mapstructure:"RUN_API_PORT"`
}

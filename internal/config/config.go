package config

type Config struct {
	KafkaBrokers        string `mapstructure:"KAFKA_BROKERS"`
	PivotDSN            string `mapstructure:"PIVOT_DSN"`
	PivotSchemaPath     string `mapstructure:"PIVOT_SCHEMA_PATH"`
	TargetDSN           string `mapstructure:"TARGET_DSN"`
	MappingPath         string `mapstructure:"MAPPING_PATH"`
	BatchMaxRows        int    `mapstructure:"BATCH_MAX_ROWS"`
	BatchMaxInterval    int    `mapstructure:"BATCH_MAX_INTERVAL"`
	JoinWaitTTL         int    `mapstructure:"JOIN_WAIT_TTL"`
	KafkaGroupID        string `mapstructure:"KAFKA_GROUP_ID"`
	KafkaCounsumerReset bool   `mapstructure:"KAFKA_CONSUMER_RESET"`
	RunApiPort          int    `mapstructure:"RUN_API_PORT"`
}

type DatabaseCredential struct {
	Type   string
	Host   string
	Port   int
	User   string
	Pass   string
	DbName string
}

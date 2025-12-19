package app

import (
	"context"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"
	"fmt"
	"strconv"
)

// LoadPlan memuat mapping root dan membuat planner.
func LoadPlan(cfg *config.Config) (*mapping.Planner, error) {
	util.Info.Println("app: loading mapping root")
	root, err := mapping.Load(*cfg)
	if err != nil {
		return nil, err
	}
	util.Info.Printf("app: mapping loaded entities=%d", len(root.Entities))

	// create planner
	return mapping.NewPlanner(root), nil
}

// InitPivot membuka repo pivot dan memastikan schema.
func InitPivot(ctx context.Context, dsn string) (*pivot.Repo, error) {
	util.Info.Println("app: connecting pivot repository")
	pivotRepo, err := pivot.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	defaultSchemaSQL := pivot.DefaultSchemaSQL()

	util.Info.Println("app: ensuring pivot schema")
	if err := pivotRepo.EnsureSchema(ctx, defaultSchemaSQL); err != nil {
		pivotRepo.Close()
		return nil, err
	}
	util.Info.Println("app: pivot schema ensured")
	return pivotRepo, nil
}

// config keys list
var DefaultConfigKey = []string{
	"TARGET_DATABASE_TYPE",
	"TARGET_DATABASE",
	"TARGET_HOST",
	"TARGET_PORT",
	"TARGET_USER",
	"TARGET_PASSWORD",
	"TARGET_DSN",
	"SOURCE_DATABASE_TYPE",
	"SOURCE_DATABASE",
	"SOURCE_HOST",
	"SOURCE_PORT",
	"SOURCE_USER",
	"SOURCE_PASSWORD",
	"DEBEZIUM_HOST",
	"DEBEZIUM_PORT",
	"DEBEZIUM_CONNECTOR_NAME",
	"PIVOT_DSN",
	"PIVOT_SCHEMA_PATH",
	"MAPPING_PATH",
	"BATCH_MAX_ROWS",
	"BATCH_MAX_INTERVAL",
	"KAFKA_BROKERS",
	"KAFKA_GROUP_ID",
	"KAFKA_CONSUMER_RESET",
	"RUN_API_PORT",
}

// get config from env and insert to pivot db then set to cfg
func ApplyBaseConfig(config *config.Config, pivotRepo *pivot.Repo) (*config.Config, error) {

	// set combined env values
	config.TargetDSN = getTargetDBConfig(config)

	values := map[string]string{
		"TARGET_DATABASE_TYPE":    config.TargetDatabaseType,
		"TARGET_DATABASE":         config.TargetDatabaseName,
		"TARGET_HOST":             config.TargetDatabaseHost,
		"TARGET_PORT":             fmt.Sprintf("%d", config.TargetDatabasePort),
		"TARGET_USER":             config.TargetDatabaseUser,
		"TARGET_PASSWORD":         config.TargetDatabasePass,
		"TARGET_DSN":              config.TargetDSN,
		"SOURCE_DATABASE_TYPE":    config.SourceDatabaseType,
		"SOURCE_DATABASE":         config.SourceDatabaseName,
		"SOURCE_HOST":             config.SourceDatabaseHost,
		"SOURCE_PORT":             fmt.Sprintf("%d", config.SourceDatabasePort),
		"SOURCE_USER":             config.SourceDatabaseUser,
		"SOURCE_PASSWORD":         config.SourceDatabasePass,
		"DEBEZIUM_HOST":           config.DebeziumHost,
		"DEBEZIUM_PORT":           fmt.Sprintf("%d", config.DebeziumPort),
		"DEBEZIUM_CONNECTOR_NAME": config.DebeziumConnectorName,
		"PIVOT_DSN":               config.PivotDSN,
		"PIVOT_SCHEMA_PATH":       config.PivotSchemaPath,
		"MAPPING_PATH":            config.MappingPath,
		"BATCH_MAX_ROWS":          fmt.Sprintf("%d", config.BatchMaxRows),
		"BATCH_MAX_INTERVAL":      fmt.Sprintf("%d", config.BatchMaxInterval),
		"KAFKA_BROKERS":           config.KafkaBrokers,
		"KAFKA_GROUP_ID":          config.KafkaGroupID,
		"KAFKA_CONSUMER_RESET":    fmt.Sprintf("%t", config.KafkaCounsumerReset),
		"RUN_API_PORT":            fmt.Sprintf("%d", config.RunApiPort),
	}

	for _, key := range DefaultConfigKey {
		val := values[key]

		cfgRow, err := pivotRepo.GetConfigurationByName(key)
		if err != nil {
			return nil, err
		}

		switch {
		case cfgRow == nil:
			if err := pivotRepo.InsertConfigurationByName(&pivot.Configuration{ConfigKey: key, ConfigValue: val}); err != nil {
				return nil, err
			}

		// update with env value
		case cfgRow.ConfigValue != val:
			if err := pivotRepo.UpdateConfigurationByName(&pivot.Configuration{ConfigKey: key, ConfigValue: val}); err != nil {
				return nil, err
			}

		default:
			val = cfgRow.ConfigValue
		}

		applyConfigValue(config, key, val)
	}

	return config, nil
}

func getTargetDBConfig(cfg *config.Config) string {

	db_module := cfg.TargetDatabaseType
	if db_module != "" && (db_module == "postgres" || db_module == "postgresql") {
		db_module = cfg.TargetDatabaseType
	} else {
		db_module = "postgres"
	}

	targetDSN := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		db_module,
		cfg.TargetDatabaseUser, cfg.TargetDatabasePass,
		cfg.TargetDatabaseHost, cfg.TargetDatabasePort,
		cfg.TargetDatabaseName,
	)

	return targetDSN
}

func applyConfigValue(cfg *config.Config, key, val string) {
	switch key {
	case "TARGET_DATABASE_TYPE":
		cfg.TargetDatabaseType = val
	case "TARGET_DATABASE":
		cfg.TargetDatabaseName = val
	case "TARGET_HOST":
		cfg.TargetDatabaseHost = val
	case "TARGET_PORT":
		cfg.TargetDatabasePort = parseInt(val, cfg.TargetDatabasePort)
	case "TARGET_USER":
		cfg.TargetDatabaseUser = val
	case "TARGET_PASSWORD":
		cfg.TargetDatabasePass = val
	case "TARGET_DSN":
		cfg.TargetDSN = val
	case "SOURCE_DATABASE_TYPE":
		cfg.SourceDatabaseType = val
	case "SOURCE_DATABASE":
		cfg.SourceDatabaseName = val
	case "SOURCE_HOST":
		cfg.SourceDatabaseHost = val
	case "SOURCE_PORT":
		cfg.SourceDatabasePort = parseInt(val, cfg.SourceDatabasePort)
	case "SOURCE_USER":
		cfg.SourceDatabaseUser = val
	case "SOURCE_PASSWORD":
		cfg.SourceDatabasePass = val
	case "DEBEZIUM_HOST":
		cfg.DebeziumHost = val
	case "DEBEZIUM_PORT":
		cfg.DebeziumPort = parseInt(val, cfg.DebeziumPort)
	case "DEBEZIUM_CONNECTOR_NAME":
		cfg.DebeziumConnectorName = val
	case "PIVOT_DSN":
		cfg.PivotDSN = val
	case "PIVOT_SCHEMA_PATH":
		cfg.PivotSchemaPath = val
	case "MAPPING_PATH":
		cfg.MappingPath = val
	case "BATCH_MAX_ROWS":
		cfg.BatchMaxRows = parseInt(val, cfg.BatchMaxRows)
	case "BATCH_MAX_INTERVAL":
		cfg.BatchMaxInterval = parseInt(val, cfg.BatchMaxInterval)
	case "KAFKA_BROKERS":
		cfg.KafkaBrokers = val
	case "KAFKA_GROUP_ID":
		cfg.KafkaGroupID = val
	case "KAFKA_CONSUMER_RESET":
		cfg.KafkaCounsumerReset = parseBool(val, cfg.KafkaCounsumerReset)
	case "RUN_API_PORT":
		cfg.RunApiPort = parseInt(val, cfg.RunApiPort)
	}
}

func parseInt(val string, fallback int) int {
	if v, err := strconv.Atoi(val); err == nil {
		return v
	}
	return fallback
}

func parseBool(val string, fallback bool) bool {
	if v, err := strconv.ParseBool(val); err == nil {
		return v
	}
	return fallback
}

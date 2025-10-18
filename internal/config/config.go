package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DefaultConfigPath string `mapstructure:"DEFAULT_CONFIG_PATH"`
	KafkaBrokers      string `mapstructure:"KAFKA_BROKERS"`
	PivotDSN          string `mapstructure:"PIVOT_DSN"`
	TargetDSN         string `mapstructure:"TARGET_DSN"`
	MappingPath       string `mapstructure:"MAPPING_PATH"`
	BatchMaxRows      int    `mapstructure:"BATCH_MAX_ROWS"`
	BatchMaxInterval  int    `mapstructure:"BATCH_MAX_INTERVAL"`
	JoinWaitTTL       int    `mapstructure:"JOIN_WAIT_TTL"`
	KafkaGroupID      string `mapstructure:"KAFKA_GROUP_ID"`
}

func Load() *Config {
	v := viper.New()

	v.SetDefault("DEFAULT_CONFIG_PATH", "./internal/config/json/default_config.json")

	v.SetConfigFile(".env")
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Tidak dapat membaca file .env: %v", err)
	}

	v.AutomaticEnv()

	// Unmarshal ke struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Gagal unmarshal config: %v", err)
	}

	fmt.Println("📦 Config Loaded:")
	fmt.Printf("Kafka brokers: %s\n", cfg.KafkaBrokers)
	fmt.Printf("Pivot DSN: %s\n", cfg.PivotDSN)
	fmt.Printf("Mapping Path: %s\n", cfg.MappingPath)

	return &cfg
}

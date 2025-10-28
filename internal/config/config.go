package config

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DefaultConfigPath   string `mapstructure:"DEFAULT_CONFIG_PATH"`
	KafkaBrokers        string `mapstructure:"KAFKA_BROKERS"`
	PivotDSN            string `mapstructure:"PIVOT_DSN"`
	TargetDSN           string `mapstructure:"TARGET_DSN"`
	MappingPath         string `mapstructure:"MAPPING_PATH"`
	BatchMaxRows        int    `mapstructure:"BATCH_MAX_ROWS"`
	BatchMaxInterval    int    `mapstructure:"BATCH_MAX_INTERVAL"`
	JoinWaitTTL         int    `mapstructure:"JOIN_WAIT_TTL"`
	KafkaGroupID        string `mapstructure:"KAFKA_GROUP_ID"`
	KafkaCounsumerReset bool   `mapstructure:"KAFKA_CONSUMER_RESET"`
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

	// reset kafka group id each run
	if cfg.KafkaCounsumerReset {
		suffix := randomString(6)
		cfg.KafkaGroupID = fmt.Sprintf("%s-%s", strings.TrimSpace(cfg.KafkaGroupID), suffix)
		fmt.Printf("🔁 Kafka group ID reset aktif. Group ID baru: %s\n", cfg.KafkaGroupID)
	}

	return &cfg
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}

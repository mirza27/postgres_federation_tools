package main

import (
	"log"

	config "db_migrate_server/config"
	replication "db_migrate_server/replication"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Load mapping configuration
	mappingConfig, err := config.LoadMappingConfig("config/ojol_mapping.json")
	if err != nil {
		log.Fatal("Error loading mapping config:", err)
	}

	
	// Load target database configuration
	targetDBConfig, err := config.LoadDBConfig("config/target.yaml")
	if err != nil {
		log.Fatal("Error loading target DB config:", err)
	}
	

	// Load source database configuration
	sourceDBConfig, err := config.LoadDBConfig("config/source.yaml")
	if err != nil {
		log.Fatal("Error loading source DB config:", err)
	}

	// Connect to target database
	targetConnStr := config.BuildConnectionString(*targetDBConfig)
	targetDB, err := sqlx.Connect("postgres", targetConnStr)
	if err != nil {
		log.Fatal("Error connecting to target DB:", err)
	}
	defer targetDB.Close()
	
	// Test target database connection
	if err := targetDB.Ping(); err != nil {
		log.Fatal("Error pinging target DB:", err)
	}
	
	log.Println("Connected to target database successfully")
	// log.Fatal("Stop here")

	// Start replication process
	if err := replication.StartReplication(*sourceDBConfig, targetDB, mappingConfig); err != nil {
		log.Fatal("Error in replication process:", err)
	}
}
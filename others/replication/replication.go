package replication

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"db_migrate_server/models"
	"db_migrate_server/processor"

	"github.com/jmoiron/sqlx"
)

// StartReplication memulai proses replikasi dengan wal2json
func StartReplication(sourceConfig models.DatabaseConfig, targetDB *sqlx.DB, mappingConfig *models.MappingConfig) error {
	// Start pg_recvlogical process dengan opsi wal2json
	cmd := exec.Command("pg_recvlogical",
		"-h", sourceConfig.Host,
		"-U", sourceConfig.User,
		"-d", sourceConfig.DBName,
		"-p 5433",
		"--slot", "db_sync_slot",
		"--start",
		"-o", "pretty-print=0",
		"-o", "add-msg-prefixes=wal2json",
		"-f", "-")

	// Set password in environment
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", sourceConfig.Password))

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting stdout pipe: %v", err)
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting pg_recvlogical: %v", err)
	}

	// Create scanner to read output
	scanner := bufio.NewScanner(stdout)

	log.Println("Starting replication process with wal2json...")

	// Process changes
	for scanner.Scan() {
		line := scanner.Text()
		var message models.Wal2JsonMessage

		fmt.Println(line)

		if err := json.Unmarshal([]byte(line), &message); err != nil {
			log.Printf("Error parsing JSON change: %v\n", err)
			continue
		}

		// Process each change item
		for _, change := range message.Change {
			if err := processor.ProcessChange(change, mappingConfig, targetDB); err != nil {
				log.Printf("Error processing change: %v\n", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from pg_recvlogical: %v", err)
	}

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("pg_recvlogical finished with error: %v", err)
	}

	return nil
}
package processor

import (
	"fmt"
	"log"

	"db_migrate_server/models"

	"github.com/jmoiron/sqlx"
)

// ProcessChange memproses perubahan dari WAL dalam format wal2json
func ProcessChange(change models.Wal2JsonChange, config *models.MappingConfig, targetDB *sqlx.DB) error {
	// Find mapping for this table
	tableMapping := findTableMapping(change.Table, config)
	if tableMapping == nil {
		return fmt.Errorf("no mapping found for table %s", change.Table)
	}

	fmt.Println("Change got :")
	fmt.Println(change)

	// Process based on operation type
	switch change.Kind {
	case "insert":
		return ProcessInsert(change, tableMapping, config, targetDB)
	case "update":
		return ProcessUpdate(change, tableMapping, config, targetDB)
	case "delete":
		return ProcessDelete(change, tableMapping, config, targetDB)
	default:
		log.Printf("Unknown operation type: %s for table %s", change.Kind, change.Table)
		return nil
	}
}

// findTableMapping mencari mapping untuk tabel tertentu
func findTableMapping(tableName string, config *models.MappingConfig) *models.TableMapping {
	for _, mapping := range config.Mappings {
		if mapping.SourceTable == tableName {
			return &mapping
		}
	}
	return nil
}

// ProcessInsert memproses operasi INSERT
func ProcessInsert(change models.Wal2JsonChange, tableMapping *models.TableMapping, config *models.MappingConfig, targetDB *sqlx.DB) error {
	for _, targetTable := range tableMapping.TargetTables {
		if !targetTable.OperationHandling.Insert {
			continue
		}

		// Build INSERT query
		query, values, err := BuildInsertQuery(change, targetTable, config)
		if err != nil {
			return fmt.Errorf("error building insert query: %v", err)
		}

		fmt.Printf("Insert Query : \n %s", query)

		// Execute query
		_, err = targetDB.Exec(query, values...)
		if err != nil {
			return fmt.Errorf("error executing insert: %v", err)
		}

		log.Printf("Inserted into %s: %v", targetTable.Table, values)
	}

	return nil
}

// ProcessUpdate memproses operasi UPDATE
func ProcessUpdate(change models.Wal2JsonChange, tableMapping *models.TableMapping, config *models.MappingConfig, targetDB *sqlx.DB) error {
	for _, targetTable := range tableMapping.TargetTables {
		if !targetTable.OperationHandling.Update {
			continue
		}

		// Build UPDATE query
		query, values, err := BuildUpdateQuery(change, targetTable, config)
		if err != nil {
			return fmt.Errorf("error building update query: %v", err)
		}

		// Execute query
		_, err = targetDB.Exec(query, values...)
		if err != nil {
			return fmt.Errorf("error executing update: %v", err)
		}

		log.Printf("Updated %s: %v", targetTable.Table, values)
	}

	return nil
}

// ProcessDelete memproses operasi DELETE
func ProcessDelete(change models.Wal2JsonChange, tableMapping *models.TableMapping, config *models.MappingConfig, targetDB *sqlx.DB) error {
	for _, targetTable := range tableMapping.TargetTables {
		// Skip if delete not handled
		if targetTable.OperationHandling.Delete == "" {
			continue
		}

		if targetTable.OperationHandling.Delete == "soft_delete" {
			// Soft delete (update a flag)
			query, values, err := BuildSoftDeleteQuery(change, targetTable, config)
			if err != nil {
				return fmt.Errorf("error building soft delete query: %v", err)
			}

			_, err = targetDB.Exec(query, values...)
			if err != nil {
				return fmt.Errorf("error executing soft delete: %v", err)
			}

			log.Printf("Soft deleted from %s: %v", targetTable.Table, values)
		} else if targetTable.OperationHandling.Delete == "hard" {
			// Hard delete (remove row)
			query, values, err := BuildHardDeleteQuery(change, targetTable, config)
			if err != nil {
				return fmt.Errorf("error building hard delete query: %v", err)
			}

			_, err = targetDB.Exec(query, values...)
			if err != nil {
				return fmt.Errorf("error executing hard delete: %v", err)
			}

			log.Printf("Hard deleted from %s: %v", targetTable.Table, values)
		}
	}

	return nil
}
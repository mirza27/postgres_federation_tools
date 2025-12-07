package pivot

import (
	"db_migrate_server/internal/config"

	"os"
)

// DefaultSchemaSQL returns the SQL that ensures the required pivot tables exist.
func DefaultSchemaSQL() string {

	defaultSchemaPath := config.Load().PivotSchemaPath
	if defaultSchemaPath != "" {
		schemaSQL, err := os.ReadFile(defaultSchemaPath)
		if err != nil {
			panic("Failed to read pivot schema file: " + err.Error())
		}

		return string(schemaSQL)

	}

	return ""
}

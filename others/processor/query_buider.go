package processor

import (
	"fmt"
	"strings"
	"time"

	"db_migrate_server/models"
)

// BuildInsertQuery membangun query INSERT untuk format wal2json
func BuildInsertQuery(change models.Wal2JsonChange, targetTable models.TargetTable, config *models.MappingConfig) (string, []interface{}, error) {
	columns := []string{}
	placeholders := []string{}
	values := []interface{}{}

	valueIndex := 1

	for _, colMapping := range targetTable.ColumnMappings {
		var value interface{}
		var err error

		if colMapping.SourceExpression != "" {
			// Skip expression for now as it requires complex evaluation
			continue
		} else if colMapping.SourceColumn != "" {
			// Find value in change data
			found := false
			for i, colName := range change.ColumnNames {
				if colName == colMapping.SourceColumn {
					value = change.ColumnValues[i]
					found = true
					break
				}
			}

			if !found {
				// Use default value if specified
				if colMapping.DefaultValue != nil {
					value = colMapping.DefaultValue
				} else {
					// Skip if no value and no default
					continue
				}
			}
		} else if colMapping.DefaultValue != nil {
			value = colMapping.DefaultValue
		} else {
			// Skip if no source and no default
			continue
		}

		// Apply transformation if specified
		if colMapping.Transformation != "" {
			value, err = ApplyTransformation(value, colMapping.Transformation, config)
			if err != nil {
				return "", nil, err
			}
		}

		columns = append(columns, colMapping.TargetColumn)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, value)
		valueIndex++
	}

	if len(columns) == 0 {
		return "", nil, fmt.Errorf("no columns to insert for table %s", targetTable.Table)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		targetTable.Table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return query, values, nil
}

// BuildUpdateQuery membangun query UPDATE untuk format wal2json
func BuildUpdateQuery(change models.Wal2JsonChange, targetTable models.TargetTable, config *models.MappingConfig) (string, []interface{}, error) {
	setClauses := []string{}
	whereClauses := []string{}
	values := []interface{}{}

	valueIndex := 1

	// Process SET clause
	for _, colMapping := range targetTable.ColumnMappings {
		if colMapping.SourceColumn == "" {
			continue
		}

		// Find value in change data
		var value interface{}
		found := false
		for i, colName := range change.ColumnNames {
			if colName == colMapping.SourceColumn {
				value = change.ColumnValues[i]
				found = true
				break
			}
		}

		if !found {
			continue
		}

		// Apply transformation if specified
		if colMapping.Transformation != "" {
			transformed, err := ApplyTransformation(value, colMapping.Transformation, config)
			if err != nil {
				return "", nil, err
			}
			value = transformed
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", colMapping.TargetColumn, valueIndex))
		values = append(values, value)
		valueIndex++
	}

	// Process WHERE clause (using old keys for update)
	if change.OldKeys != nil {
		for i, keyName := range change.OldKeys.KeyNames {
			// Find the key mapping
			for _, colMapping := range targetTable.ColumnMappings {
				if colMapping.SourceColumn == keyName {
					value := change.OldKeys.KeyValues[i]

					// Apply transformation if specified
					if colMapping.Transformation != "" {
						transformed, err := ApplyTransformation(value, colMapping.Transformation, config)
						if err != nil {
							return "", nil, err
						}
						value = transformed
					}

					whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", colMapping.TargetColumn, valueIndex))
					values = append(values, value)
					valueIndex++
					break
				}
			}
		}
	}

	if len(whereClauses) == 0 {
		return "", nil, fmt.Errorf("no valid WHERE clause for update on table %s", targetTable.Table)
	}

	if len(setClauses) == 0 {
		return "", nil, fmt.Errorf("no columns to update for table %s", targetTable.Table)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		targetTable.Table,
		strings.Join(setClauses, ", "),
		strings.Join(whereClauses, " AND "))

	return query, values, nil
}

// BuildSoftDeleteQuery membangun query untuk soft delete dengan format wal2json
func BuildSoftDeleteQuery(change models.Wal2JsonChange, targetTable models.TargetTable, config *models.MappingConfig) (string, []interface{}, error) {
	setClauses := []string{}
	whereClauses := []string{}
	values := []interface{}{}

	valueIndex := 1

	// Set deletion flag
	setClauses = append(setClauses, "is_deleted = true")

	// Add timestamp if configured
	if config.ConflictResolution.TimestampColumn != "" {
		setClauses = append(setClauses,
			fmt.Sprintf("%s = $%d", config.ConflictResolution.TimestampColumn, valueIndex))
		values = append(values, time.Now())
		valueIndex++
	}

	// Process WHERE clause (using old keys)
	if change.OldKeys != nil {
		for i, keyName := range change.OldKeys.KeyNames {
			// Find the key mapping
			for _, colMapping := range targetTable.ColumnMappings {
				if colMapping.SourceColumn == keyName {
					value := change.OldKeys.KeyValues[i]

					// Apply transformation if specified
					if colMapping.Transformation != "" {
						transformed, err := ApplyTransformation(value, colMapping.Transformation, config)
						if err != nil {
							return "", nil, err
						}
						value = transformed
					}

					whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", colMapping.TargetColumn, valueIndex))
					values = append(values, value)
					valueIndex++
					break
				}
			}
		}
	}

	if len(whereClauses) == 0 {
		return "", nil, fmt.Errorf("no valid WHERE clause for delete on table %s", targetTable.Table)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		targetTable.Table,
		strings.Join(setClauses, ", "),
		strings.Join(whereClauses, " AND "))

	return query, values, nil
}

// BuildHardDeleteQuery membangun query untuk hard delete dengan format wal2json
func BuildHardDeleteQuery(change models.Wal2JsonChange, targetTable models.TargetTable, config *models.MappingConfig) (string, []interface{}, error) {
	whereClauses := []string{}
	values := []interface{}{}

	valueIndex := 1

	// Process WHERE clause (using old keys)
	if change.OldKeys != nil {
		for i, keyName := range change.OldKeys.KeyNames {
			// Find the key mapping
			for _, colMapping := range targetTable.ColumnMappings {
				if colMapping.SourceColumn == keyName {
					value := change.OldKeys.KeyValues[i]

					// Apply transformation if specified
					if colMapping.Transformation != "" {
						transformed, err := ApplyTransformation(value, colMapping.Transformation, config)
						if err != nil {
							return "", nil, err
						}
						value = transformed
					}

					whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", colMapping.TargetColumn, valueIndex))
					values = append(values, value)
					valueIndex++
					break
				}
			}
		}
	}

	if len(whereClauses) == 0 {
		return "", nil, fmt.Errorf("no valid WHERE clause for delete on table %s", targetTable.Table)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s",
		targetTable.Table,
		strings.Join(whereClauses, " AND "))

	return query, values, nil
}
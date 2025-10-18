package processor

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"db_migrate_server/models"
)

// ApplyTransformation menerapkan transformasi pada nilai
func ApplyTransformation(value interface{}, transformation string, config *models.MappingConfig) (interface{}, error) {
	transformationFunc, exists := config.Transformations[transformation]
	if !exists {
		return nil, fmt.Errorf("transformation %s not found", transformation)
	}

	// Simple transformation handling
	switch transformation {
	case "string_to_int":
		if str, ok := value.(string); ok {
			var intVal int
			_, err := fmt.Sscanf(str, "%d", &intVal)
			if err != nil {
				return nil, err
			}
			return intVal, nil
		}
		return value, nil
	case "date_to_string":
		if t, ok := value.(time.Time); ok {
			return t.Format("2006-01-02"), nil
		}
		if s, ok := value.(string); ok {
			// Try to parse as date
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t.Format("2006-01-02"), nil
			}
		}
		return value, nil
	case "json_to_text":
		if s, ok := value.(string); ok {
			return s, nil
		}
		// For non-string values, convert to JSON string
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		return string(jsonBytes), nil
	default:
		// Handle custom transformations defined in config
		if strings.Contains(transformationFunc, "{value}") {
			result := strings.ReplaceAll(transformationFunc, "{value}", fmt.Sprintf("%v", value))
			return result, nil
		}
		return value, nil
	}
}
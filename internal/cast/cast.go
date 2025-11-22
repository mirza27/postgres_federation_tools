package cast

import (
	"fmt"
	"strconv"
	"strings"
)

// Value applies a simple cast instruction to the provided value, returning the
// converted result or an error when the requested type is unsupported.
func Value(cast string, value interface{}) (interface{}, error) {
	if cast == "" || value == nil {
		return value, nil
	}

	switch strings.ToLower(cast) {
	case "string", "text", "varchar":
		return fmt.Sprint(value), nil
	case "bool", "boolean":
		return toBool(value)
	case "int", "integer", "bigint", "smallint":
		return toInt64(value)
	case "numeric", "decimal", "float", "double", "float8", "real":
		return toFloat64(value)
	default:
		return nil, fmt.Errorf("unsupported cast %s", cast)
	}
}

func toBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int8:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case uint:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return parseBoolString(v)
	default:
		return parseBoolString(fmt.Sprint(value))
	}
}

func parseBoolString(s string) (bool, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return false, nil
	}
	if b, err := strconv.ParseBool(s); err == nil {
		return b, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f != 0, nil
	}
	return false, fmt.Errorf("cannot cast %q to bool", s)
}

func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return parseIntString(v)
	default:
		return parseIntString(fmt.Sprint(value))
	}
}

func parseIntString(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(f), nil
	}
	return 0, fmt.Errorf("cannot cast %q to int", s)
}

func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return parseFloatString(v)
	default:
		return parseFloatString(fmt.Sprint(value))
	}
}

func parseFloatString(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	return 0, fmt.Errorf("cannot cast %q to float", s)
}

package expr

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Evaluate attempts to resolve a known expression string into a Go value. The
// boolean indicates whether the expression was handled; callers can fall back
// to other logic when it is false.
func Evaluate(raw string) (interface{}, bool, error) {
	name, args := parse(raw)
	if name == "" {
		return nil, false, nil
	}

	switch name {
	case "now":
		return time.Now().UTC(), true, nil
	case "uuid", "uuid_v4":
		return uuid.NewString(), true, nil
	case "random_string":
		length := 16
		if len(args) >= 1 && args[0] != "" {
			val, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, true, fmt.Errorf("random_string invalid length: %w", err)
			}
			length = val
		}
		str, err := randomString(length)
		return str, true, err
	case "random_number", "random_int":
		min, max, err := parseNumberRange(args)
		if err != nil {
			return nil, true, err
		}
		val, err := randomInt(min, max)
		return val, true, err
	default:
		return nil, false, nil
	}
}

func parse(raw string) (string, []string) {
	expr := strings.TrimSpace(raw)
	if expr == "" {
		return "", nil
	}
	open := strings.Index(expr, "(")
	close := strings.LastIndex(expr, ")")
	if open == -1 || close == -1 || close < open {
		return strings.ToLower(expr), nil
	}
	name := strings.ToLower(strings.TrimSpace(expr[:open]))
	argsStr := strings.TrimSpace(expr[open+1 : close])
	if argsStr == "" {
		return name, nil
	}
	parts := strings.Split(argsStr, ",")
	var args []string
	for _, p := range parts {
		args = append(args, strings.TrimSpace(p))
	}
	return name, args
}

func randomString(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		buf[i] = alphabet[idx.Int64()]
	}
	return string(buf), nil
}

func parseNumberRange(args []string) (int64, int64, error) {
	var (
		min int64 = 0
		max int64 = 999999
	)
	if len(args) == 1 {
		val, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("random_number invalid max: %w", err)
		}
		max = val
	} else if len(args) >= 2 {
		valMin, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("random_number invalid min: %w", err)
		}
		valMax, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("random_number invalid max: %w", err)
		}
		min, max = valMin, valMax
	}
	if max < min {
		min, max = max, min
	}
	return min, max, nil
}

func randomInt(min, max int64) (int64, error) {
	if max <= min {
		return min, nil
	}
	rangeSize := max - min + 1
	limit := big.NewInt(rangeSize)
	val, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return 0, err
	}
	return min + val.Int64(), nil
}

package format

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

func CountDigits(n int) int {
	var count int = 0

	if n == 0 {
		return 1
	}
	if n < 0 {
		n = -n
	}

	for n > 0 {
		n /= 10
		count++
	}

	return count
}

func GetValueAsString(val any) string {
	switch val := val.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return strconv.FormatBool(val)
	case time.Time:
		return val.Format("2006-01-02 15:04:05")
	case map[string]any, []any:
		b, err := json.Marshal(val)
		if err != nil {
			return "<invalid json>"
		}
		return string(b)
	}
	slog.Error("Unhandled type while formatting", slog.Any("type", fmt.Sprintf("%T", val)))
	return fmt.Sprintf("%v", val)
}

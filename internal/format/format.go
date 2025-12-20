package format

import (
	"fmt"
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
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", val)
	case time.Time:
		// @TODO: Consider specific formatting
		return fmt.Sprintf("%v", val)
	}
	return fmt.Sprintf("ERR: UNHANDLED TYPE '%T'", val)
}

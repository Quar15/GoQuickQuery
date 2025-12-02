package utilities

import (
	"fmt"
	"time"
)

const KeySmallG int = 103
const KeySmallH int = 104
const KeySmallJ int = 106
const KeySmallK int = 107
const KeySmallL int = 108
const KeySmallV int = 118

func Contains(slice []int, element int) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

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

func DebugPrintMap(m []map[string]any) {
	for i, data := range m {
		for k, v := range data {
			fmt.Printf("idx = %v | key = %v | val =  %v | type=%T\n", i, k, v, v)
		}
	}
}

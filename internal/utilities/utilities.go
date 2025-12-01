package utilities

import (
	"strconv"
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
	case int:
		return strconv.Itoa(int(val))
	}
	return ""
}

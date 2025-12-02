package utilities

import (
	"strconv"
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
	case int:
		return strconv.Itoa(int(val))
	}
	return ""
}

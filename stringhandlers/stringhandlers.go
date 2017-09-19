package stringhandlers

import (
	"strconv"
	"strings"
	"time"
)

// ConvertField converts the values from string to the type configured
func ConvertField(typeData, format, value string) (newValue interface{}, err error) {

	switch typeData {
	case "date":
		newValue, err = time.Parse(format, strings.TrimSpace(value))
		break
	case "number":
		newValue, err = strconv.ParseFloat(strings.TrimSpace(value), 64)
		break
	default:
		newValue = value
	}
	return
}

// Substr extracts the chars from string using runes
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length

	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}

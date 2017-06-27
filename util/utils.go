package util

import (
	"strconv"
	"strings"
	"time"
)

// ConvertField - Convert values from string to the type configured
func ConvertField(typeData, value string) (newValue interface{}, err error) {

	switch typeData {
	case "date":
		newValue, err = time.Parse(time.RFC3339, strings.TrimSpace(value))
		break
	case "number":
		newValue, err = strconv.ParseFloat(strings.TrimSpace(value), 64)
		break
	default:
		newValue = value
	}
	return
}

// Substr - Extract chars from string using runes
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length

	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}

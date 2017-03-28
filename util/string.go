package util

import (
	"strings"
)

// ContainsString searches the provided string to see if it contains *any* of the
// provided search terms. The caseSensitive parameter is used to determine if
// the matches should be case sensitive.
func ContainsString(str string, terms []string, caseSensitive bool) bool {
	if !caseSensitive {
		str = strings.ToLower(str)
		termsLower := make([]string, len(terms))
		for i, _ := range terms {
			termsLower[i] = strings.ToLower(terms[i])
		}
		terms = termsLower
	}
	for _, t := range terms {
		if strings.Contains(str, t) {
			return true
		}
	}
	return false
}

// SplitAndTrimString splits a string using the provided separator and then
// trims spaces from each of the strings.
func SplitAndTrimString(str, separator string) []string {
	slice := strings.Split(str, separator)
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}
	return slice
}

// Truncate ensures that a message consists of less than 140 characters,
// appending an ellipsis if necessary.
func Truncate(message string, maxLen int) string {
	var v = message
	if len(v) > maxLen {
		v = strings.TrimSpace(v[:maxLen-1])
		v += "…"
	}
	return v
}

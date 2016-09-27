package util

import (
	"strings"
)

// ContainsString searches the provided string to see if it contains *any* of the
// provided search terms. The caseSensitive parameter is used to determine if
// the matches should be case sensitive.
func ContainsString(terms []string, str string, caseSensitive bool) bool {
	if !caseSensitive {
		str = strings.ToLower(str)
		termsLower := make([]string, len(terms))
		for i, _ := range terms {
			termsLower[i] = strings.ToLower(terms[i])
		}
		terms = termsLower
	}
	for _, t := range terms {
		if strings.Contains(t, str) {
			return true
		}
	}
	return false
}

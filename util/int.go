package util

import (
	"strconv"
)

// Atoi converts a string into an integer, returning 0 upon error.
func Atoi(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

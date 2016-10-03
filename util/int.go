package util

import (
	"strconv"
)

// Atoi converts a string into an integer, returning 0 upon error.
func Atoi(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

// Atoi64 converts a string into a 64-bit integer, returning 0 upon error.
func Atoi64(str string) int64 {
	v, _ := strconv.ParseInt(str, 10, 64)
	return v
}

// ContainsInt determines whether the integer slice contains the specified
// integer.
func ContainsInt(integers []int, v int) bool {
	for _, i := range integers {
		if i == v {
			return true
		}
	}
	return false
}

// FilterInt executes the provided function for each integer in the slice and
// only includes it in the slice returned if the function returns true.
func FilterInt(integers []int, fn func(int) bool) (newIntegers []int) {
	for _, i := range integers {
		if fn(i) {
			newIntegers = append(newIntegers, i)
		}
	}
	return
}

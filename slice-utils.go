package validator

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// SortSlice sorts a slice of any orderable value
func SortSlice[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

// RemoveDuplicatesInSortedSlice removes duplicates from a sorted slice
func RemoveDuplicatesInSortedSlice[T constraints.Ordered](s []T) []T {
	if len(s) < 1 {
		return s
	}

	n := 1
	for i := 1; i < len(s); i++ {
		if s[i-1] != s[i] {
			s[n] = s[i]
			n++
		}
	}

	return s[:n]
}

package utils

import "github.com/zijiren233/gencontainer/restrictions"

func EqualSlice[T restrictions.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

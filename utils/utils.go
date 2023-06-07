package utils

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](left T, right T) T {
	if left < right {
		return left
	} else {
		return right
	}
}

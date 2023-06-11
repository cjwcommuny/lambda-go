package num

import "golang.org/x/exp/constraints"

func Increment[T constraints.Integer](x T) T {
	return x + 1
}

func IncrementBy[T constraints.Integer](by T) func(T) T {
	return func(x T) T {
		return x + by
	}
}

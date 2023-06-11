package collection

import "github.com/cjwcommuny/lambda-go/adt/opt"

func MapGetter[K comparable, V any](m map[K]V) func(K) opt.Option[V] {
	return func(key K) opt.Option[V] {
		value, ok := m[key]
		if ok {
			return opt.Some(value)
		} else {
			return opt.None[V]()
		}
	}
}

func SliceIndex[T any](slice []T) func(int) T {
	return func(i int) T {
		return slice[i]
	}
}

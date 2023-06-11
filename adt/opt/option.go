package opt

import "github.com/barweiss/go-tuple"

type Option[T any] struct {
	some    bool
	payload T
}

func Some[T any](x T) Option[T] {
	return Option[T]{
		some:    true,
		payload: x,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		some: false,
	}
}

func IsSome[T any](o Option[T]) bool {
	return o.some
}

func GetSomeUnchecked[T any](o Option[T]) T {
	return o.payload
}

func IsNone[T any](o Option[T]) bool {
	return !o.some
}

func Map[A any, B any](f func(A) B) func(Option[A]) Option[B] {
	return func(o Option[A]) Option[B] {
		if IsSome(o) {
			return Some(f(o.payload))
		} else {
			return None[B]()
		}
	}
}

func Zip[A any, B any](optionA Option[A]) func(Option[B]) Option[tuple.T2[A, B]] {
	return func(optionB Option[B]) Option[tuple.T2[A, B]] {
		if IsSome(optionA) && IsSome(optionB) {
			return Some(tuple.New2(GetSomeUnchecked(optionA), GetSomeUnchecked(optionB)))
		} else {
			return None[tuple.T2[A, B]]()
		}
	}
}

func UnwrapOr[T any](defaultValue T) func(Option[T]) T {
	return func(o Option[T]) T {
		if IsSome(o) {
			return GetSomeUnchecked(o)
		} else {
			return defaultValue
		}
	}
}

func UnwrapOrElse[T any](defaultFunc func() T) func(Option[T]) T {
	return func(o Option[T]) T {
		if IsSome(o) {
			return GetSomeUnchecked(o)
		} else {
			return defaultFunc()
		}
	}
}

func Expect[T any](panicValue any) func(Option[T]) T {
	return func(o Option[T]) T {
		if IsSome(o) {
			return GetSomeUnchecked(o)
		} else {
			panic(panicValue)
		}
	}
}

func Unwrap[T any](o Option[T]) T {
	if IsSome(o) {
		return GetSomeUnchecked(o)
	} else {
		panic("called `Unwrap` on a `None` value")
	}
}

func Inspect[T any](f func(T)) func(Option[T]) {
	return func(o Option[T]) {
		if IsSome(o) {
			f(GetSomeUnchecked(o))
		}
	}
}

func Or[T any](left Option[T]) func(Option[T]) Option[T] {
	return func(right Option[T]) Option[T] {
		if IsSome(left) {
			return left
		} else {
			return right
		}
	}
}

func And[T any](left Option[T]) func(Option[T]) Option[T] {
	return func(right Option[T]) Option[T] {
		if IsNone(left) {
			return None[T]()
		} else {
			return right
		}
	}
}

func Flatten[T any](o Option[Option[T]]) Option[T] {
	if IsSome(o) {
		return GetSomeUnchecked(o)
	} else {
		return None[T]()
	}
}

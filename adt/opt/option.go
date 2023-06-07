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

func GetSome[T any](o Option[T]) T {
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
			return Some(tuple.New2(GetSome(optionA), GetSome(optionB)))
		} else {
			return None[tuple.T2[A, B]]()
		}
	}
}

package iter

import (
	"github.com/barweiss/go-tuple"
	"github.com/cjwcommuny/lambda-go/adt"
	"github.com/cjwcommuny/lambda-go/adt/opt"
	"github.com/cjwcommuny/lambda-go/num"
	"github.com/cjwcommuny/lambda-go/utils"
)

type SizeHint struct {
	LowerBound int
	UpperBound opt.Option[int]
}

type Iter[E any] struct {
	next     func() opt.Option[E]
	sizeHint SizeHint
}

func Next[E any](iter Iter[E]) opt.Option[E] {
	return iter.next()
}

func GetSizeHint[E any](iter Iter[E]) SizeHint {
	return iter.sizeHint
}

func NewIterWithoutSizeHint[E any](next func() opt.Option[E]) Iter[E] {
	return Iter[E]{
		next: next,
		sizeHint: SizeHint{
			LowerBound: 0,
			UpperBound: opt.None[int](),
		},
	}
}

func NewIterWithStaticSizeHintLowerBound[E any](next func() opt.Option[E], sizeHint int) Iter[E] {
	return NewIter(
		next,
		SizeHint{
			LowerBound: sizeHint,
			UpperBound: opt.Some(sizeHint),
		},
	)
}

func NewIter[E any](next func() opt.Option[E], sizeHint SizeHint) Iter[E] {
	return Iter[E]{next, sizeHint}
}

func SliceIter[E any](slice []E) Iter[E] {
	index := 0
	return Iter[E]{
		next: func() opt.Option[E] {
			if index < len(slice) {
				element := slice[index]
				index++
				return opt.Some(element)
			} else {
				return opt.None[E]()
			}
		},
		sizeHint: SizeHint{
			LowerBound: len(slice),
			UpperBound: opt.Some(len(slice)),
		},
	}
}

func Map[A any, B any](f func(A) B) func(Iter[A]) Iter[B] {
	return func(base Iter[A]) Iter[B] {
		return Iter[B]{
			next: func() opt.Option[B] {
				element := base.next()
				return opt.Map(f)(element)
			},
			sizeHint: GetSizeHint(base),
		}
	}
}

func ForEach[E any](f func(E)) func(Iter[E]) adt.Void {
	return func(it Iter[E]) adt.Void {
		for {
			element := it.next()
			if opt.IsNone(element) {
				break
			}
			f(opt.GetSomeUnchecked(element))
		}
		return adt.MakeVoid()
	}
}

func Find[E any](predicate func(E) bool) func(Iter[E]) opt.Option[E] {
	return func(it Iter[E]) opt.Option[E] {
		for {
			element := it.next()
			if opt.IsNone(element) {
				return opt.None[E]()
			}
			if predicate(opt.GetSomeUnchecked(element)) {
				return element
			}
		}
	}
}

func Filter[E any](predicate func(E) bool) func(Iter[E]) Iter[E] {
	return func(it Iter[E]) Iter[E] {
		return Iter[E]{
			next: func() opt.Option[E] {
				return Find(predicate)(it)
			},
			sizeHint: SizeHint{
				LowerBound: 0,
				UpperBound: GetSizeHint(it).UpperBound,
			},
		}
	}
}

func Fold[B any, E any](init B, f func(B, E) B) func(Iter[E]) B {
	return func(it Iter[E]) B {
		ForEach(func(e E) {
			init = f(init, e)
		})(it)
		return init
	}
}

func Zip[A any, B any](iterA Iter[A]) func(Iter[B]) Iter[tuple.T2[A, B]] {
	return func(iterB Iter[B]) Iter[tuple.T2[A, B]] {
		sizeHintA := GetSizeHint(iterA)
		sizeHintB := GetSizeHint(iterB)
		upperBound := func() opt.Option[int] {
			if opt.IsSome(sizeHintA.UpperBound) && opt.IsSome(sizeHintB.UpperBound) {
				return opt.Some(utils.Min(opt.GetSomeUnchecked(sizeHintA.UpperBound), opt.GetSomeUnchecked(sizeHintB.UpperBound)))
			} else if opt.IsSome(sizeHintA.UpperBound) {
				return sizeHintA.UpperBound
			} else if opt.IsSome(sizeHintB.UpperBound) {
				return sizeHintB.UpperBound
			} else {
				return opt.None[int]()
			}
		}()
		next := func() opt.Option[tuple.T2[A, B]] {
			elementA := iterA.next()
			elementB := iterB.next()
			return opt.Zip[A, B](elementA)(elementB)
		}
		return NewIter(
			next,
			SizeHint{
				LowerBound: utils.Min(sizeHintA.LowerBound, sizeHintA.LowerBound),
				UpperBound: upperBound,
			},
		)
	}
}

func CollectToSlice[E any](it Iter[E]) []E {
	result := make([]E, 0, GetSizeHint(it).LowerBound)
	ForEach(func(e E) {
		result = append(result, e)
	})(it)
	return result
}

func CollectToMap[K comparable, V any](it Iter[tuple.T2[K, V]]) map[K]V {
	result := make(map[K]V, GetSizeHint(it).LowerBound)
	ForEach(func(t tuple.T2[K, V]) {
		result[t.V1] = t.V2
	})(it)
	return result
}

func Count[E any](it Iter[E]) int {
	counter := 0
	ForEach(func(element E) {
		counter++
	})(it)
	return counter
}

func Iterate[E any](generate func(E) E) func(init E) Iter[E] {
	return func(init E) Iter[E] {
		previous := opt.None[E]()
		next := func() opt.Option[E] {
			return opt.MapOrElse(
				func() opt.Option[E] {
					previous = opt.Some(init)
					return opt.Some(init)
				},
				func(pre E) opt.Option[E] {
					current := generate(pre)
					previous = opt.Some(current)
					return opt.Some(current)
				},
			)(previous)
		}
		sizeHint := SizeHint{LowerBound: 1, UpperBound: opt.None[int]()}
		return NewIter(next, sizeHint)
	}
}

func Take[E any](n int) func(Iter[E]) Iter[E] {
	return func(iterator Iter[E]) Iter[E] {
		count := 0
		next := func() opt.Option[E] {
			if count < n {
				count++
				return iterator.next()
			} else {
				return opt.None[E]()
			}
		}
		originalSizeHint := GetSizeHint(iterator)
		sizeHint := SizeHint{
			LowerBound: originalSizeHint.LowerBound,
			UpperBound: opt.Some(opt.MapOr(
				n,
				func(upperBound int) int { return utils.Min(n, upperBound) },
			)(originalSizeHint.UpperBound)),
		}
		return NewIter(next, sizeHint)
	}
}

func Enumerate[E any](it Iter[E]) Iter[tuple.T2[int, E]] {
	indices := Iterate(num.Increment[int])(0)
	return Zip[int, E](indices)(it)
}

func TakeWhile[E any](predicate func(E) bool) func(Iter[E]) Iter[E] {
	return func(iterator Iter[E]) Iter[E] {
		next := func() opt.Option[E] {
			element := iterator.next()
			return opt.AndThen(func(e E) opt.Option[E] {
				if predicate(e) {
					return opt.Some(e)
				} else {
					return opt.None[E]()
				}
			})(element)
		}
		return NewIter(
			next,
			SizeHint{LowerBound: 0, UpperBound: GetSizeHint(iterator).UpperBound},
		)
	}
}

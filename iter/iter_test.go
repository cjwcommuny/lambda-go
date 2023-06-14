package iter

import (
	"github.com/barweiss/go-tuple"
	"github.com/cjwcommuny/lambda-go/adt/opt"
	"github.com/cjwcommuny/lambda-go/num"
	"math"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/cjwcommuny/lambda-go/fn"
)

func quickTest(t *testing.T, f any) {
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestForEach(t *testing.T) {
	f := func(s []int) bool {
		result := make([]int, 0)
		fn.Pipe2(
			SliceIter[int],
			ForEach(func(element int) { result = append(result, element) }),
		)(s)
		return reflect.DeepEqual(s, result)
	}
	quickTest(t, f)
}

func TestMap(t *testing.T) {
	mapFunc := func(x int) int { return x + 1 }
	f := func(s []int) bool {
		result := fn.Pipe3(
			SliceIter[int],
			Map(mapFunc),
			CollectToSlice[int],
		)(s)
		sameLength := len(result) == len(s)
		sameElements := func() bool {
			for index := range s {
				if mapFunc(s[index]) != result[index] {
					return false
				}
			}
			return true
		}()
		return sameElements && sameLength
	}
	quickTest(t, f)
}

func TestCount(t *testing.T) {
	f := func(s []int) bool {
		count := fn.Pipe2(SliceIter[int], Count[int])(s)
		return count == len(s)
	}
	quickTest(t, f)
}

func TestZip(t *testing.T) {
	f := func(s1 []int, s2 []int) bool {
		zipped := fn.Pipe3(
			SliceIter[int],
			Zip[int, int](SliceIter(s1)),
			CollectToSlice[tuple.T2[int, int]],
		)(s2)
		expectedLength := int(math.Min(float64(len(s1)), float64(len(s2))))
		sameLength := len(zipped) == expectedLength
		sameElements := func() bool {
			for index := 0; index < expectedLength; index++ {
				if s1[index] != zipped[index].V1 || s2[index] != zipped[index].V2 {
					return false
				}
			}
			return true
		}()
		return sameLength && sameElements
	}
	quickTest(t, f)
}

func TestFold(t *testing.T) {
	f := func(s []int) bool {
		result := fn.Pipe2(
			SliceIter[int],
			Fold(make([]int, 0), func(slice []int, element int) []int {
				return append(slice, element)
			}),
		)(s)
		return reflect.DeepEqual(s, result)
	}
	quickTest(t, f)
}

func TestFilter(t *testing.T) {
	f := func(s []int) bool {
		result := fn.Pipe3(
			SliceIter[int],
			Filter(func(x int) bool { return x > 0 }),
			CollectToSlice[int],
		)(s)
		return func() bool {
			count := 0
			for _, x := range s {
				if x > 0 {
					if result[count] != x {
						return false
					}
					count++
				}
			}
			return true
		}()
	}
	quickTest(t, f)
}

func TestFind(t *testing.T) {
	f := func(s []int) bool {
		result := fn.Pipe2(
			SliceIter[int],
			Find(func(x int) bool { return x > 0 }),
		)(s)
		expectedResult := func() opt.Option[int] {
			for _, x := range s {
				if x > 0 {
					return opt.Some(x)
				}
			}
			return opt.None[int]()
		}()
		return reflect.DeepEqual(result, expectedResult)
	}
	quickTest(t, f)
}

func TestIterateAndTake(t *testing.T) {
	f := func(init uint, step uint, count uint) bool {
		init = init % 100
		step = step % 100
		count = count % 20
		result := fn.Pipe3(
			Iterate(num.IncrementBy(step)),
			Take[uint](int(count)),
			CollectToSlice[uint],
		)(init)
		lengthMatch := len(result) == int(count)
		elementsMatch := func() bool {
			for index := 0; index < int(count); index++ {
				if result[index] != init+uint(index)*step {
					return false
				}
			}
			return true
		}()
		return lengthMatch && elementsMatch
	}
	quickTest(t, f)
}

func TestEnumerate(t *testing.T) {
	f := func(s []int) bool {
		result := fn.Pipe3(
			SliceIter[int],
			Enumerate[int],
			CollectToSlice[tuple.T2[int, int]],
		)(s)
		for index := range s {
			if result[index].V1 != index || result[index].V2 != s[index] {
				return false
			}
		}
		return true
	}
	quickTest(t, f)
}

func TestSkip(t *testing.T) {
	f := func(s []int, n uint16) bool {
		skip := func() int {
			if len(s) == 0 {
				return int(n)
			} else {
				return int(n) % len(s)
			}
		}()
		result := fn.Pipe3(
			SliceIter[int],
			Skip[int](skip),
			CollectToSlice[int],
		)(s)
		lengthMatch := len(result) == int(math.Max(float64(len(s)-skip), 0))
		elementsMatch := func() bool {
			for index := skip; index < len(s); index++ {
				if s[index] != result[index-skip] {
					return false
				}
			}
			return true
		}()
		return lengthMatch && elementsMatch
	}
	quickTest(t, f)
}

func TestSkipWhile(t *testing.T) {
	f := func(s []int) bool {
		result := fn.Pipe3(
			SliceIter[int],
			SkipWhile(func(x int) bool { return x > 0 }),
			CollectToSlice[int],
		)(s)
		expectedResult := make([]int, 0)
		index := 0
		for ; index < len(s); index++ {
			if s[index] <= 0 {
				break
			}
		}
		for ; index < len(s); index++ {
			expectedResult = append(expectedResult, s[index])
		}
		return reflect.DeepEqual(result, expectedResult)
	}
	quickTest(t, f)
}

func TestReduce(t *testing.T) {
	f := func(s []int) bool {
		result := fn.Pipe2(
			SliceIter[int],
			Reduce(func(left int, right int) int { return left + right }),
		)(s)
		expectedResult := opt.None[int]()
		for _, x := range s {
			expectedResult = opt.Some(opt.MapOr(x, func(acc int) int { return acc + x })(expectedResult))
		}
		return reflect.DeepEqual(result, expectedResult)
	}
	quickTest(t, f)
}

func TestTakeWhile(t *testing.T) {
	f := func(s []int) bool {
		result := fn.Pipe3(
			SliceIter[int],
			TakeWhile(func(x int) bool { return x > 0 }),
			CollectToSlice[int],
		)(s)
		expectedResult := make([]int, 0)
		for _, x := range s {
			if x > 0 {
				expectedResult = append(expectedResult, x)
			} else {
				break
			}
		}
		return reflect.DeepEqual(result, expectedResult)
	}
	quickTest(t, f)
}

func TestCollectToMap(t *testing.T) {
	f := func(s []int) bool {
		result := fn.Pipe3(
			SliceIter[int],
			Map(func(x int) tuple.T2[int, int] { return tuple.New2(x, x/2) }),
			CollectToMap[int, int],
		)(s)
		expectedResult := make(map[int]int)
		for _, x := range s {
			expectedResult[x] = x / 2
		}
		return reflect.DeepEqual(result, expectedResult)
	}
	quickTest(t, f)
}

package iter

import (
	"github.com/barweiss/go-tuple"
	"github.com/cjwcommuny/lambda-go/adt/opt"
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

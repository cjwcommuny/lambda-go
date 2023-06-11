package iter

import (
	"github.com/barweiss/go-tuple"
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

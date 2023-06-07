package iter

import (
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

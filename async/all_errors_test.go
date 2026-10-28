package async_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/zeiss/pkg/async"
)

func TestAllPropagatesErrors(t *testing.T) {
	want := errors.New("worker failed")
	for _, failed := range [][]bool{{true, false}, {false, true}, {true, true}} {
		t.Run(fmt.Sprint(failed), func(t *testing.T) {
			futures := make([]*async.Future[int], len(failed))
			for i, fail := range failed {
				futures[i] = async.New(func() (int, error) {
					if fail {
						return 0, fmt.Errorf("task %d: %w", i, want)
					}
					return i + 1, nil
				})
			}
			combined := async.All(futures...)
			_, got := combined.Await()
			if !errors.Is(got, want) {
				t.Fatalf("Await error = %v; want wrapped worker error", got)
			}
			calls := 0
			combined.Catch(func(err error) {
				calls++
				if !errors.Is(err, want) {
					t.Errorf("Catch error = %v; want wrapped worker error", err)
				}
			})
			if calls != 1 {
				t.Fatalf("Catch called %d times; want 1", calls)
			}
		})
	}
}

func TestAllSuccessfulValues(t *testing.T) {
	for _, want := range [][]int{{}, {10, 20, 30}} {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			futures := make([]*async.Future[int], len(want))
			for i, value := range want {
				futures[i] = async.New(func() (int, error) { return value, nil })
			}
			got, err := async.All(futures...).Await()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %v; want %v", got, want)
			}
		})
	}
}

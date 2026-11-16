package channels

import (
	"sort"
	"testing"
	"time"
)

func TestFanInEmptyInputs(t *testing.T) {
	for name, fanIn := range map[string]func(...<-chan int) <-chan int{"Join": Join[int], "Merge": Merge[int]} {
		t.Run(name, func(t *testing.T) {
			select {
			case _, ok := <-fanIn():
				if ok {
					t.Fatal("expected a closed output")
				}
			case <-time.After(time.Second):
				t.Fatal("output did not close for zero inputs")
			}
		})
	}
}

func TestFanInWaitsForRemainingInput(t *testing.T) {
	for name, fanIn := range map[string]func(...<-chan int) <-chan int{"Join": Join[int], "Merge": Merge[int]} {
		t.Run(name, func(t *testing.T) {
			empty := make(chan int)
			close(empty)
			remaining := make(chan int, 1)
			defer close(remaining)
			out := fanIn(empty, remaining)
			select {
			case value, ok := <-out:
				t.Fatalf("output finished before the remaining input: value=%v, open=%v", value, ok)
			case <-time.After(100 * time.Millisecond):
			}
			remaining <- 42
			select {
			case value, ok := <-out:
				if !ok || value != 42 {
					t.Fatalf("got %v, %v; want 42, true", value, ok)
				}
			case <-time.After(time.Second):
				t.Fatal("remaining value was not forwarded")
			}
		})
	}
}

func TestFanInDrainsEveryInput(t *testing.T) {
	for name, fanIn := range map[string]func(...<-chan int) <-chan int{"Join": Join[int], "Merge": Merge[int]} {
		t.Run(name, func(t *testing.T) {
			inputs := make([]<-chan int, 3)
			for i := range inputs {
				input := make(chan int, 2)
				input <- i * 2
				input <- i*2 + 1
				close(input)
				inputs[i] = input
			}
			out := fanIn(inputs...)
			got := []int{}
			timeout := time.After(time.Second)
			for {
				select {
				case value, ok := <-out:
					if !ok {
						sort.Ints(got)
						if len(got) != 6 {
							t.Fatalf("got %v; want every input value", got)
						}
						for i, value := range got {
							if i != value {
								t.Fatalf("got %v; want 0 through 5", got)
							}
						}
						return
					}
					got = append(got, value)
				case <-timeout:
					t.Fatal("output did not close after all inputs were drained")
				}
			}
		})
	}
}

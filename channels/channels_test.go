package channels

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	in1 := make(chan int)
	in2 := make(chan int)

	out := Join(in1, in2)

	go func() {
		in1 <- 1
		in2 <- 2
	}()

	el := []int{}

	assert.Eventually(t, func() bool {
		e := <-out
		el = append(el, e)
		return len(el) == 2
	}, 1*time.Second, 10*time.Millisecond)
}

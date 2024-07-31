package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zeiss/pkg/async"
)

func main() {
	fn := func() (string, error) {
		timer := time.NewTimer(3 * time.Second)
		defer timer.Stop()

		<-timer.C

		return "hello", nil
	}

	errFn := func() (string, error) {
		return "", fmt.Errorf("this is a main error")
	}

	async.New(fn).Then(func(s string) {
		log.Print(s)
	})

	async.New(errFn).Then(func(s string) {
		log.Print(s)
	}).Catch(func(err error) {
		log.Print(fmt.Errorf("error: %w", err))
	})
}

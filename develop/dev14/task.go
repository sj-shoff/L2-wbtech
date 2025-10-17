package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	start := time.Now()

	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		cases := make([]reflect.SelectCase, len(channels))
		for i, ch := range channels {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}

		reflect.Select(cases)
	}()

	return orDone
}

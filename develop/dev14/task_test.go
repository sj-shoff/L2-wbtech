package main

import (
	"testing"
	"time"
)

func TestOrDone(t *testing.T) {
	fast := make(chan interface{})
	mid := make(chan interface{})
	slow := make(chan interface{})

	go func() {
		time.Sleep(50 * time.Millisecond)
		close(fast)
	}()

	go func() {
		time.Sleep(75 * time.Millisecond)
		close(mid)
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		close(slow)
	}()

	start := time.Now()
	<-or(fast, mid, slow)
	elapsed := time.Since(start)

	if elapsed < 45*time.Millisecond || elapsed > 55*time.Millisecond {
		t.Error(
			"expected", 50*time.Millisecond,
			"got", elapsed,
		)
	}

}

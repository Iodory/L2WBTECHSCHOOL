package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	if len(channels) == 0 {
		close(out)
		return out
	}

	if len(channels) == 1 {
		return channels[0]
	}

	go func() {
		defer close(out)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			mid := len(channels) / 2
			select {
			case <-or(channels[:mid]...):
			case <-or(channels[mid:]...):
			}
		}
	}()
	return out
}

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
	fmt.Printf("done after %v\n", time.Since(start))
}

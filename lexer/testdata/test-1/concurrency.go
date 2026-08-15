package main

import "fmt"

func concurrencyDemo() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()

	ch := make(chan int, 1)
	done := make(chan struct{})

	go func() {
		ch <- 42
		close(done)
	}()

	select {
	case v := <-ch:
		fmt.Println("got:", v)
	case <-done:
	}

	<-done

	panic("ludicrous panic for demonstration")
}

package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	MaxOutstanding = 10
)

type Request struct {
	count int
}

var sem = make(chan int, MaxOutstanding)

func process(r *Request) {
	fmt.Printf("wake up, count=%d\n", r.count)
	time.Sleep( time.Duration(r.count*10) * time.Millisecond)
}

/*
func handle(r *Request) {
	fmt.Printf("before process\n ")
    sem <- 1    // Wait for active queue to drain.
    process(r)  // May take a long time.
    <-sem       // Done; enable next request to run.
}*/

func Serve(queue chan *Request) {
	var wg sync.WaitGroup

	for r := range queue {
		fmt.Printf("before process, count=%d\n ", r.count)
		sem <- 1
		go func(r *Request) {
			wg.Add(1)
			process(r)
			
			<-sem
			wg.Done()
		}(r)

		//req := <-queue
		//go handle(req)  // Don't wait for handle to finish.
	}

	fmt.Printf("no more to process, break\n ")
	wg.Wait()
}
func produce(queue chan *Request) {
	for i := 0; i < 100; i++ {
		queue <- &Request{i}
		time.Sleep(10 * time.Microsecond)
	}
	close(queue)

}
func main() {
	queue := make(chan *Request, 15)
	go produce(queue)
	Serve(queue)
}

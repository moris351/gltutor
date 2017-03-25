package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
)

func fibonacci(c chan int64, quit chan int) {
	x, y := int64(0), int64(1)
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int64)
	quit := make(chan int)
	go func() {
		for i := 0; i < 77; i++ {
			fmt.Printf("%d: %d\n", i, <-c)
			time.Sleep(100*time.Millisecond)
		}
		//quit <- 0
	}()
	go fibonacci(c, quit)

	fmt.Println("Hello world!")
  	fmt.Print("Press 'Enter' to exit...\n")
  	bufio.NewReader(os.Stdin).ReadBytes('\n') 
	quit<-0

}
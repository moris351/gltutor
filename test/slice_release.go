package main

import (
	"bytes"
	"fmt"
	_ "io/ioutil"
	"runtime"
	"strings"
	"time"
)

type messager interface {
	id()
}
type message struct {
	msg string
	len int
}

func (m *message) id() {
	fmt.Printf("msg=%s len=%d\n", m.msg, m.len)
}

func main() {
	//sl := []messagerr{{"anna",1}, {"crist",1}}
	in := make(chan messager)
	out := make(chan messager)
	go consume(out)
	go sliceIQ(in, out)
	go pump(in)

	//time.Sleep(100*time.Second)
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		printMemStats()
	}

}
func pump(in chan messager) {
	for {
		fmt.Println("pump ")
		in <- &message{"don", 1}
	}

}

func consume(out chan messager) {
	for {
		fmt.Println("consume ")
		m := <-out
		m.id()
	}

}

func sliceIQ(in <-chan messager, next chan<- messager) {
	defer close(next)

	// pending events (this is the "infinite" part)
	pending := []messager{}

recv:
	for {
		// Ensure that pending always has values so the select can
		// multiplex between the receiver and sender properly
		fmt.Println("sliceIQ 1")
		if len(pending) == 0 {
			v, ok := <-in

			if !ok {
				// in is closed, flush values
				fmt.Println("in is closed, flush values")
				break
			}
			fmt.Println("sliceIQ 2")

			// We now have something to send
			pending = append(pending, v)
		}
		fmt.Println("sliceIQ 3")

		select {
		// Queue incoming values
		case v, ok := <-in:
			if !ok {
				// in is closed, flush values
				break recv
			}
			pending = append(pending, v)
			fmt.Println("sliceIQ 4")

		// Send queued values
		case next <- pending[0]:
			pending[0] = nil
			pending = pending[1:]
			fmt.Println("sliceIQ 5")
		}
	}

	// After in is closed, we may still have events to send
	for _, v := range pending {
		next <- v
		fmt.Println("sliceIQ 6")
	}
}

/*
func pipe(sl []messager, msg messager)  []messager {
	sll:=[]messager{}
	sll=append(sll,msg)
	sl = append(sl, msg)
	sll[0] = nil
	sl = sl[1:]
	return sl
}*/

func printMemStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	buf := bytes.NewBuffer(nil)
	buf.WriteString(strings.Repeat("=", 72) + "\n")
	buf.WriteString("Memory Profile:\n")
	buf.WriteString(fmt.Sprintf("\tAlloc: %d Kb\n", mem.Alloc/1024))
	buf.WriteString(fmt.Sprintf("\tTotalAlloc: %d Kb\n", mem.TotalAlloc/1024))
	buf.WriteString(fmt.Sprintf("\tNumGC: %d\n", mem.NumGC))
	buf.WriteString(fmt.Sprintf("\tGoroutines: %d\n", runtime.NumGoroutine()))
	buf.WriteString(strings.Repeat("=", 72))
	fmt.Println(buf.String())
}

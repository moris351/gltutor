package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"gltutor/webcrawler/fetcher"
	"io"
	"os"
	"runtime"
	_ "strconv"
	"strings"
	"sync"
	"time"

	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
)

var (
	memStats = flag.Duration("memstats", 0, "display memory statistics at a given interval")
)

//go:generate gotemplate "github.com/ncw/gotemplate/set" mySet(string)

func main() {
	flag.Parse()
	// Pass a log message and arguments directly
	logger := log.Logger()
	//appender:=logger.Appender()

	logger.SetAppender(appenders.RollingFile("webcrawler.log", true))
	//logger.SetAppender(appenders.Console())
	appender := logger.Appender()
	//alayout := appender.Layout()
	appender.SetLayout(layout.Pattern("%d-%p-%m"))

	inputFile, inputError := os.Open("input.dat")

	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n")
		return // exit the function on error
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	var inputString []string

	for {
		str, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {
			break
		}
		inputString = append(inputString, strings.Trim(str, "\015\012"))
		//aq = append(aq, make(chan string))
	}
	for _, str := range inputString {
		log.Debug("The input was: %s\n", str)
	}
	//cont := make(chan string)

	f := fetcher.New()

	// First mem stat print must be right after creating the fetchbot

	log.Debug("*memStats=%v", *memStats)
	if *memStats > 0 {
		// Print starting stats
		log.Debug("*memStats=%v", *memStats)
		printMemStats(nil)
		// Run at regular intervals
		runMemStats(f, *memStats)
		// On exit, print ending stats after a GC
		defer func() {
			runtime.GC()
			printMemStats(nil)
		}()
	}

	f.Start(inputString)

	//output := bufio.NewWriter(outputFile)
	/*
		for _, req := range inputString {
			str := <-cont
			output.WriteString(req)
			output.WriteString(str)
			fmt.Printf("%s \n", str)

		}
	*/
	//fmt.Printf("%s", robots)
	fmt.Println("Shutting Down")
}

func runMemStats(f *fetcher.Fetcher, tick time.Duration) {
	var mu sync.Mutex
	var di *fetcher.DebugInfo

	// Start goroutine to collect fetchbot debug info
	go func() {
		for v := range f.Debug() {
			mu.Lock()
			di = v
			mu.Unlock()
		}
	}()
	// Start ticker goroutine to print mem stats at regular intervals
	go func() {
		c := time.Tick(tick)
		for _ = range c {
			mu.Lock()
			printMemStats(di)
			mu.Unlock()
		}
	}()
}

func printMemStats(di *fetcher.DebugInfo) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	buf := bytes.NewBuffer(nil)
	buf.WriteString(strings.Repeat("=", 72) + "\n")
	buf.WriteString("Memory Profile:\n")
	buf.WriteString(fmt.Sprintf("\tAlloc: %d Kb\n", mem.Alloc/1024))
	buf.WriteString(fmt.Sprintf("\tTotalAlloc: %d Kb\n", mem.TotalAlloc/1024))
	buf.WriteString(fmt.Sprintf("\tNumGC: %d\n", mem.NumGC))
	buf.WriteString(fmt.Sprintf("\tGoroutines: %d\n", runtime.NumGoroutine()))
	if di != nil {
		buf.WriteString(fmt.Sprintf("\tNumHosts: %d\n", di.NumHosts))
	}
	buf.WriteString(strings.Repeat("=", 72))
	log.Debug(buf.String())
}

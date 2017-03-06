package main

import (
	"bufio"
	"fmt"
	"gltutor/webcrawler/fetcher"
	"io"
	"os"
	_ "runtime"
	_ "strconv"
	"strings"
	_ "time"

	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/log"
)

func main() {
	// Pass a log message and arguments directly
	logger := log.Logger()
	//appender:=logger.Appender()

	logger.SetAppender(appenders.RollingFile("webcrawler.log", true))
	//logger.SetAppender(appenders.Console())
	appender:=logger.Appender()
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
	for _,str:=range inputString{
		log.Debug("The input was: %s\n", str)
	}
	//cont := make(chan string)

	f := fetcher.New()
	f.Start(inputString)

	outputFile, outputError := os.Create("output.dat")

	if outputError != nil {
		fmt.Printf("An error occurred on opening the outputfile\n")
		return // exit the function on error
	}
	defer outputFile.Close()

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

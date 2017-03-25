package main

import (
    "os"
    "bufio"
    "fmt"
)

func main () {
    // var outputWriter *bufio.Writer
    // var outputFile *os.File
    // var outputError os.Error
    // var outputString string
	dup:=make(map[string]bool)
	dup["item1"]=true
	dup["item2"]=false
	if flag,ok:=dup["item1"];!ok{
		fmt.Printf("item1 dup,flag=%b",flag)
	}
	if flag,ok:=dup["item3"];!ok{
		fmt.Printf("item3 dup,flag=%b",flag)
	}
	
    outputFile, outputError := os.OpenFile("output.dat", os.O_APPEND|os.O_CREATE, 0666)
    if outputError != nil {
        fmt.Printf("An error occurred with file opening or creation\n")
        return  
    }
    defer outputFile.Close()

    outputWriter := bufio.NewWriter(outputFile)
    outputString := "hello world!\n"

    for i:=0; i<10; i++ {
        outputWriter.WriteString(outputString)
    }
    outputWriter.Flush()
}
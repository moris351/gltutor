package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	_ "runtime"
	_ "strconv"
	"strings"
	_ "time"
	"gltutor/webcrawer/fetcher"

	"github.com/PuerkitoBio/goquery"
	"github.com/ian-kent/go-log/log"
)

func main() {
	// Pass a log message and arguments directly
	log.Debug("Example log message: %s", "example arg")

	inputFile, inputError := os.Open("input.dat")

	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return // exit the function on error
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	i := 0
	var inputString []string
	//var aq []chan string

	//var readerError error

	for {
		str, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {
			break
		}
		inputString = append(inputString, strings.Trim(str, "\015\012"))
		//aq = append(aq, make(chan string))
		fmt.Printf("The input was:%d %s\n", i, inputString[i])
		i++
	}
	cont := make(chan string)
	/*for i, url := range inputString {
		if len(url) == 0 {
			break
		}

		//log.Debug("url=%s",url)
		//go fetch(url, cont,aq[i])
	}*/
	f := fetcher.New()
	f.Start(inputString)

	outputFile, outputError := os.Create("output.dat")

	if outputError != nil {
		fmt.Printf("An error occurred on opening the outputfile\n")
		return // exit the function on error
	}
	defer outputFile.Close()

	output := bufio.NewWriter(outputFile)

	for _, req := range inputString {
		str := <-cont
		output.WriteString(req)
		output.WriteString(str)
		fmt.Printf("%s \n", str)

	}
	//fmt.Printf("%s", robots)
	fmt.Println("Shutting Down")
}
func fetch(url string, cont chan string, aqi chan string ) (err error) {

	fmt.Printf("url=%s\n", url)
	//url = "http://www.google.com/robots.txt"
	resp, err := http.Get(strings.Trim(url, "\015\012"))
	if err != nil {
		log.Fatal("error message: %s", err)
		return err
	}

	defer resp.Body.Close()
	log.Debug("Readall!")
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		//fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
		return
	}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		aqi <- val 

	})

	c, err := ioutil.ReadAll(resp.Body)
	//str := string(c[:2])
	cont <- string(c[:])
	//res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}

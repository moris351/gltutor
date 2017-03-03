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
	//var readerError error

	for {
		str, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {
			break
		}
		inputString = append(inputString, str)
		fmt.Printf("The input was:%d %s", i, inputString[i])
		i++
	}
	/*var inputString [3]string
	inputString[0]="http://www.google.com/robots.txt"
	inputString[1]="http://www.google.com/robots.txt"
	inputString[2]="http://www.google.com/robots.txt"
	*/
	// Pass a function which returns a log message and arguments
	//log.Debug(func() { []interface{}{"Example log message: %s", "example arg"} })
	//log.Debug(func(i ...interface{}) { []interface{}{"Example log message: %s", "example arg"} })
	//var cont chan string
	cont := make(chan string)
	for _, url := range inputString {
		if len(url) == 0 {
			break
		}

		//log.Debug("url=%s",url)
		go fetch(url, cont)
	}
	for range inputString {
		input := <-cont
		fmt.Printf("%s \n", input)
	}
	//fmt.Printf("%s", robots)
	fmt.Println("Shutting Down")
}
func fetch(url string, cont chan string) (err error) {

	fmt.Printf("url=%s\n", url)
	//url = "http://www.google.com/robots.txt"
	resp, err := http.Get(strings.Trim(url, "\015\012"))
	if err != nil {
		log.Fatal("error message: %s", err)
		return err
	}

	defer resp.Body.Close()
	log.Debug("Readall!")

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

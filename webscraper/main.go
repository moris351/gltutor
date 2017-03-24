package main

import(
	"fmt"
	"io/ioutil"
	"net/http"
	_ "time"
)
func main() {
	client := http.Client{}/*
		Timeout: time.Duration(15 * time.Second),
	}*/
	resp, err := client.Get("http://www.github.com")
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("error message: %s\n", err)
		return 
	}
	content,err:=ioutil.ReadAll(resp.Body)

	fmt.Printf("resp:%s\n",content)
	return 
}

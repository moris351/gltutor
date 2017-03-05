package fetcher

import (
	_ "bufio"
	_ "fmt"
	_ "io"
	_ "io/ioutil"
	"net/http"
	"net/url"
	_ "os"
	_ "runtime"
	_ "strconv"
	_ "strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ian-kent/go-log/log"
)

type Command struct{
	u url.URL
	p string
}
//fetcher struct
type Fetcher struct {
	urls map[string]chan Command
	//chUrl []chan url.URL
	delay time.Duration
}

func New() *Fetcher{
	return &Fetcher{
		delay : 5 * time.Second,
		urls : make(map[string]chan Command),
	}
}

func (f *Fetcher) getUrls(host string) (chan Command){
		_,err:=f.urls[host]
		if err==false {
			f.urls[host]=make(chan Command,10)
		}
		return f.urls[host]
}
//pop url from chan
//call handle
//loop for next
func (f *Fetcher) doRequest()  {
	log.Debug("doRequest")
	for _, cmd := range f.urls{
		go f.parseChan(cmd)
		//f.handle(cmd)
	}
}

func (f *Fetcher)parseChan(cmd chan Command){
	for {
		select {
        	case v := <-cmd:
				log.Debug("Received on cmd: %s\n", v.u.String())
				go f.handle(v)
		}
		time.Sleep(f.delay)
	}
}

//visit url
//fetch html
//find a anchor
//add a anchor to chan
func (f *Fetcher) handle(cmd Command) error {
	log.Debug("handle: ru=%s",cmd.u.String())
	resp, err := http.Get(cmd.u.String())
	if err != nil {
		log.Fatal("error message: %s", err)
		return err
	}

	defer resp.Body.Close()
	log.Debug("Readall!")
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		//fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
		return err
	}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		log.Debug("doc.find: val=",val)
		u, err := cmd.u.Parse(val)
		if err!=nil {
			log.Debug("parse failed")
			return
		}
		log.Debug("u=%s",u.String())
		f.getUrls(u.Host)<-Command{*u,""}
 

	})
	return err
}

//add original url to chan
//call doRequest
func (f *Fetcher) Start(rawUrls []string)  {
	log.Debug("start")
		
	for _, ru := range rawUrls {
		log.Debug("start ru:%s",ru)
		u,err:=url.Parse(ru)
		if err!=nil {
			log.Debug("url.Parse failed!")
			break;
		}
		
		ch := f.getUrls(u.Host)

		ch<-Command{*u,""}
	}

	f.doRequest()
	
}
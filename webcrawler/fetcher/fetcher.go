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
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ian-kent/go-log/log"
)

type Command struct {
	u url.URL
	p string
}

type PageInfo struct{
	u url.URL
	title string
	description string
}

//fetcher struct
type Fetcher struct {
	urls map[string]chan Command
	back chan PageInfo
	dup map[string]bool

	deep map[string]int
	deepLimit int
	//chUrl []chan url.URL
	delay time.Duration
}

func New() *Fetcher {
	return &Fetcher{
		delay: 5 * time.Second,
		urls:  make(map[string]chan Command),
		back: make(chan PageInfo),
		deep: make(map[string]int),
		deepLimit: 3,
	}
}
/*
func (f *Fetcher) getUrls(host string) chan Command {
	ch, err := f.urls[host]
	if err == false {
		f.urls[host] = make(chan Command, 10)
	}
	return f.urls[host]
}
*/
//pop url from chan
//call handle
//loop for next
func (f *Fetcher) doRequest() {
	log.Debug("doRequest")
	for _, cmd := range f.urls {
		go f.parseChan(cmd)
		//f.handle(cmd)
	}
}

func (f *Fetcher) parseBack(back chan PageInfo){
	log.Debug("parseBack")
	for {
		info,ok:=<-back
		if(!ok){log.Debug("info read not ok")}
		log.Debug("u=%s", info.u.String())
		ss:= info.u.String()
		if i:=strings.Index(ss,"#"); i !=-1{
			ss = ss[0:i]
		}
		log.Debug("ss=%s", ss)
		
		if _,ok:=f.dup[ss];!ok{
			if _,ok:=f.urls[info.u.Host];!ok {
				log.Debug("new host=%s",info.u.Host)
				f.urls[info.u.Host] = make(chan Command,3)
				go f.parseChan(f.urls[info.u.Host])
			}
			f.urls[info.u.Host] <- Command{info.u, ""}
		}else
		{
			log.Debug("url dup, url=%s",info.u.String())
			f.dup[ss]=true
		}

	}
}

func (f *Fetcher) parseChan(cmd chan Command) {
	log.Debug("parseChan")
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
	log.Debug("handle: ru=%s", cmd.u.String())
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

	title :=doc.Find("html title").Text()
	log.Debug("title=%s",title)

	description,_ :=doc.Find("head meta[name='description']").Attr("content")
	log.Debug("description=%s",description)
	/*for host,deep:=range f.deep{
		log.Debug("crew deep =%d, host=%s",deep,host)

	}
	
	if(f.deep[cmd.u.Host]>=f.deepLimit){
		log.Debug("crew deep reached, host=%s",cmd.u.Host)
		return nil
	}*/
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		log.Debug("doc.find: val=%s", val)
		u, err := cmd.u.Parse(val)
		if err != nil {
			log.Debug("parse failed")
			return
		}
		f.back<-PageInfo{*u,title,description}

	})

	//log.Debug("now deep=%d, host=%s",f.deep[cmd.u.Host],cmd.u.Host)

	//f.deep[cmd.u.Host]+=1
	return err
}

//add original url to chan
//call doRequest
func (f *Fetcher) Start(rawUrls []string) {
	log.Debug("start")

	for _, ru := range rawUrls {
		log.Debug("start ru:%s", ru)
		u, err := url.Parse(ru)
		if err != nil {
			log.Debug("url.Parse failed!")
			break
		}

		_,ok := f.urls[u.Host]
		if !ok {
			f.urls[u.Host] = make(chan Command,3)
		}
		
		go f.parseChan(f.urls[u.Host])

		f.urls[u.Host] <- Command{*u, ""}
	}
	f.back = make(chan PageInfo,3)
	f.parseBack(f.back)
	//f.doRequest()

}

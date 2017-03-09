package fetcher

import (
	"bufio"
	_ "fmt"
	_ "io"
	_ "io/ioutil"
	"net/http"
	"net/url"
	"os"
	_ "runtime"
	_ "strconv"
	"strings"
	"time"
	"sync"
	"github.com/PuerkitoBio/goquery"
	"github.com/ian-kent/go-log/log"
)
const (
	DefaultWorkerIdleTTL = 30 * time.Second
	DefaultWorkerDelay = 5 * time.Second
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
// The DebugInfo holds information to introspect the Fetcher's state.
type DebugInfo struct {
	NumHosts int
}

//fetcher struct
type Fetcher struct {
	mu    sync.Mutex
	urls map[string]chan Command
	back chan PageInfo
	dup map[string]bool

	deep map[string]int
	deepLimit int
	//chUrl []chan url.URL
	delay time.Duration
	
	WorkerIdleTTL time.Duration
	// dbg is a channel used to push debug information.
	dbgmu     sync.Mutex
	dbg       chan *DebugInfo
	debugging bool
}

func New() *Fetcher {
	return &Fetcher{
		delay: DefaultWorkerDelay,
		urls:  make(map[string]chan Command),
		back: make(chan PageInfo),
		dup: make(map[string]bool),
		deep: make(map[string]int),
		dbg: make(chan *DebugInfo),
		WorkerIdleTTL: DefaultWorkerIdleTTL,
		deepLimit: 3,
	}
}

// Debug returns the channel to use to receive the debugging information. It is not intended
// to be used by package users.
func (f *Fetcher) Debug() <-chan *DebugInfo {
	f.dbgmu.Lock()
	defer f.dbgmu.Unlock()
	f.debugging = true
	return f.dbg
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
		if f.dup == nil{
			log.Debug("f.dup is nil")
		}
		if _,ok:=f.dup[ss];!ok{
			f.dup[ss]=true
			if _,ok:=f.urls[info.u.Host];!ok {
				log.Debug("new host=%s",info.u.Host)
				f.mu.Lock()
				f.urls[info.u.Host] = make(chan Command,3)
				f.mu.Unlock()
				go f.parseChan(f.urls[info.u.Host],info.u.Host)
			}
			f.urls[info.u.Host] <- Command{info.u, ""}
		}else
		{
			log.Debug("url dup, url=%s",info.u.String())
			//f.dup[ss]=true
		}

		f.output(info)

	}
}

func (f *Fetcher) parseChan(cmd chan Command, hostkey string) {
	log.Debug("parseChan")
	
	var ttl   <-chan time.Time
	//var delay <-chan time.Time
	delay := time.NewTicker(5*time.Second)
	defer delay.Stop()
	var cmdQ []*Command

	endloop := false
	for !endloop {
		select {
		case v := <-cmd:
			log.Debug("Received on cmd: %s\n", v.u.String())
			cmdQ = append(cmdQ,&v)
			ttl = time.After(f.WorkerIdleTTL)

		// Send queued values
		case <-delay.C:
			log.Debug("ready for new search, current queue length: %d\n", len(cmdQ))
			if len(cmdQ) == 0 {break}
			if len(cmdQ) >= 1{
				v := cmdQ[0]
				cmdQ[0] = nil
				cmdQ = cmdQ[1:]
				f.handle(*v)
			}

		case <-ttl:
			// Worker has been idle for WorkerIdleTTL, terminate it
			f.mu.Lock()
			inch, ok := f.urls[hostkey]
			delete(f.urls, hostkey)
			f.mu.Unlock()
			if ok {
				close(inch)
			}
			endloop = true
			break
		}
		
		f.dbgmu.Lock()
		log.Debug("debugging=%t",f.debugging)
		if f.debugging {
			//f.mu.Lock()
			select {
			case f.dbg <- &DebugInfo{len(f.urls)}:
			default:
			}
			//f.mu.Unlock()
		}
		f.dbgmu.Unlock()
		
		//time.Sleep(f.delay)
	}
}

//visit url
//fetch html
//find a anchor
//add a anchor to chan
func (f *Fetcher) handle(cmd Command) error {
	log.Debug("handle: ru=%s", cmd.u.String())
	client := http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
	resp, err := client.Get(cmd.u.String())
	//resp, err := http.Get(cmd.u.String())
	if err != nil {
		log.Debug("error message: %s", err)
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
		
		go f.parseChan(f.urls[u.Host],u.Host)

		f.urls[u.Host] <- Command{*u, ""}
	}
	f.back = make(chan PageInfo,3)
	//f.doRequest()
	f.parseBack(f.back)
}
func (f *Fetcher) output(info PageInfo){
    outputFile, outputError := os.OpenFile("output.dat", os.O_APPEND|os.O_CREATE, 0666)
    if outputError != nil {
        log.Debug("An error occurred with file opening or creation\n")
        return  
    }
    defer outputFile.Close()

    outputWriter := bufio.NewWriter(outputFile)
    //outputString := "hello world!\n"
	outputWriter.WriteString(info.u.String()+"\n"+info.title+"\n"+info.description+"\n")

    outputWriter.Flush()


}
/*
func detectContentCharset(r reader) string {
    if data, err := r.Peek(1024); err == nil {
        if _, name, ok := charset.DetermineEncoding(data, ""); ok {
            return name
        }
    }
    return "utf8"
}*/
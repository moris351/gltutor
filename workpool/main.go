package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/goinggo/workpool"
	"github.com/ian-kent/go-log/log"
)

type MyWork struct {
	Name      string "The Name of a person"
	BirthYear int    "The Yea the person was born"
	WP        *workpool.WorkPool
}

func (workPool *MyWork) DoWork(workRoutine int) {
	fmt.Printf("%s : %d\n", workPool.Name, workPool.BirthYear)
	fmt.Printf("*******> WR: %d  QW: %d  AR: %d\n", workRoutine, workPool.WP.QueuedWork(), workPool.WP.ActiveRoutines())
	time.Sleep(100 * time.Millisecond)

	//panic("test")
}

func main() {
	// Pass a log message and arguments directly
	log.Debug("Example log message: %s", "example arg")
	
	// Pass a function which returns a log message and arguments
	//log.Debug(func() { []interface{}{"Example log message: %s", "example arg"} })
	//log.Debug(func(i ...interface{}) { []interface{}{"Example log message: %s", "example arg"} })
	runtime.GOMAXPROCS(runtime.NumCPU())

	workPool := workpool.New(runtime.NumCPU()*3, 10)

	shutdown := false // Just for testing, I Know

	go func() {

		for i := 0; i < 1000; i++ {

			work := &MyWork{
				Name:      "A" + strconv.Itoa(i),
				BirthYear: i,
				WP:        workPool,
			}

			err := workPool.PostWork("name_routine", work)

			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				time.Sleep(100 * time.Millisecond)
			}

			if shutdown == true {
				return
			}
		}
	}()

	fmt.Println("Hit any key to exit")

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	shutdown = true

	fmt.Println("Shutting Down\n")
	workPool.Shutdown("name_routine")
}

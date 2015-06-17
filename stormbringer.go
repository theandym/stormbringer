package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	)

func main() {
	// Constants
	const MaxInt = int(^uint(0) >> 1)

	// Flags
	workersFlag := flag.Int("workers", 8, "number of concurrent worker processes")
	lengthFlag := flag.Int("length", 10000, "number of calls per endpoint per worker")
	flag.Parse()

	// Workers are the number of concurrent processes used cURL target URLs
	workers := workersFlag

	// Length is the number of calls per target endpoint per worker
	var length int64
	if *lengthFlag == 0 {
		length = int64(MaxInt)
	} else {
		length = int64(*lengthFlag)
	}

	// Targets are the target URLs to cURL
	var targetsArg string
	var targets []string
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	} else {
		targetsArg = flag.Args()[0]
		st := strings.Split(targetsArg, ",")
		for _, v := range st {
			targets = append(targets, strings.TrimSpace(v))
		}
	}


	// Review
	fmt.Println("workers=" + strconv.Itoa(*workers), "length=" + strconv.FormatInt(length, 10), "targets=\"" + targetsArg + "\"")

	// Spawn `workers` goroutines
	// http://www.goinggo.net/2014/01/concurrency-goroutines-and-gomaxprocs.html
	runtime.GOMAXPROCS(*workers)
  var wg sync.WaitGroup
  for worker := 1; worker <= *workers; worker ++ {
    wg.Add(1)
    go func(worker int) {
			loadGen(worker, length, targets)
      wg.Done()
    }(worker)
  }
	wg.Wait()

	time.Sleep(time.Hour * 24 * 365)

}

// Iterate through `targets` for `length` cycles
func loadGen(worker int, length int64, targets []string) {

	var iteration int64
	for iteration = 1; iteration <= length; iteration ++ {

		rand.Seed(time.Now().UnixNano())
    shuffle(targets)

		for _, value := range targets {

			cmd := exec.Command(
				"curl",
				"-sSLw",
				"status=%{http_code} total_time=%{time_total} time_connect=%{time_connect} time_start=%{time_starttransfer} target=\"%{url_effective}\"\n",
				value,
				"-o",
				"/dev/null",
				)
			out, err := cmd.Output()

			if err != nil {
			fmt.Println(err.Error())
			return
			}

			fmt.Print(string(out))

		}

	}

}

// Shuffle slice elements.
// http://marcelom.github.io/2013/06/07/goshuffle.html
func shuffle(a []string) {
    for i := range a {
        j := rand.Intn(i + 1)
        a[i], a[j] = a[j], a[i]
    }
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
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
	curlFlag := flag.Bool("curl", false, "switch to `curl` for requests (from the Go `net/http` package)")
	workersFlag := flag.Int("workers", 4, "number of concurrent worker processes")
	lengthFlag := flag.Int("length", 10000, "number of requests per target per worker")
	flag.Parse()

	// Using the --curl flag switches to `curl` for requests (from the Go `net/http` package)
	curl := *curlFlag

	// Workers are the number of concurrent processes used to request target URLs
	workers := workersFlag

	// Length is the number of requests per target per worker
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
	fmt.Println("stormbringer=start workers="+strconv.Itoa(*workers), "length="+strconv.FormatInt(length, 10), "targets=\""+targetsArg+"\"")

	// Spawn `workers` goroutines
	// http://www.goinggo.net/2014/01/concurrency-goroutines-and-gomaxprocs.html
	runtime.GOMAXPROCS(*workers)
	var wg sync.WaitGroup
	for worker := 1; worker <= *workers; worker++ {
		wg.Add(1)
		go func(worker int) {
			loadGen(curl, worker, length, targets)
			wg.Done()
		}(worker)
	}
	wg.Wait()

	fmt.Println("stormbringer=end")
	time.Sleep(time.Hour * 24 * 365)

}

// Iterate through `targets` for `length` cycles
func loadGen(curl bool, worker int, length int64, targets []string) {

	var iteration int64
	for iteration = 1; iteration <= length; iteration++ {

		rand.Seed(time.Now().UnixNano())
		shuffle(targets)

		for _, value := range targets {

			if curl == true {

				// curl version
				cmd := exec.Command(
					"curl",
					"-sSLw",
					"worker="+strconv.Itoa(worker)+" iteration="+strconv.FormatInt(iteration, 10)+" target=\"%{url_effective}\" status=%{http_code} total_time=%{time_total} time_connect=%{time_connect} time_start=%{time_starttransfer}\n",
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

			} else {

				// Go `net/http` version
				start_time := time.Now()

				resp, err := http.Get(value)

				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				} else {
					_, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}

					end_time := time.Now()

					// fmt.Println(string(body))
					fmt.Println(
						"worker="+strconv.Itoa(worker),
						"iteration="+strconv.FormatInt(iteration, 10),
						"target=\""+value+"\"",
						"status="+strconv.Itoa(resp.StatusCode),
						"total_time="+strconv.FormatFloat(timer(start_time, end_time), 'f', 3, 64),
					)

				}
				resp.Body.Close()
			}

		}

	}

}

// Shuffle slice elements
// http://marcelom.github.io/2013/06/07/goshuffle.html
func shuffle(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// Timer
func timer(start time.Time, end time.Time) float64 {
	// time.Since(start)
	elapsed := end.Sub(start)
	return toFixed(elapsed.Seconds(), 3)
}

// `float64` truncation
// http://stackoverflow.com/a/29786394
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

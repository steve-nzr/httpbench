package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/gammazero/workerpool"
)

func createClient(wc int) *http.Client {
	return &http.Client{
		Timeout: 1 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   50 * time.Millisecond,
				KeepAlive: 10 * time.Millisecond,
			}).DialContext,
			DisableKeepAlives:   false,
			IdleConnTimeout:     2 * time.Hour,
			MaxIdleConns:        wc,
			MaxIdleConnsPerHost: wc,
			MaxConnsPerHost:     wc,
		},
	}
}

func main() {
	poolSize := flag.Int("poolsize", 100, "poolsize")
	runTimes := flag.Int("runtimes", 10, "runtimes")
	calls := flag.Int("calls", 1000, "calls")
	flag.Parse()

	runtime.GOMAXPROCS(1)

	timeouts := 0
	times := []time.Duration{}
	for i := 0; i < *runTimes; i++ {
		c := createClient(*poolSize)
		pool := workerpool.New(*poolSize)

		begin := time.Now()

		// compute
		for i := 0; i < *calls; i++ {
			pool.Submit(func() {
				res, err := c.Get("http://localhost:8080/test")
				if res != nil {
					res.Body.Close()
				}
				if err != nil {
					timeouts++
				}
			})
		}

		pool.StopWait()

		// end
		duration := time.Now().Sub(begin)
		times = append(times, duration)
		fmt.Println(fmt.Sprintf("Test serie number %d/%d finished in %dms", i+1, *runTimes, duration.Milliseconds()))
	}

	var timesAll float64
	for _, t := range times {
		timesAll += (float64)(t.Milliseconds())
	}

	fmt.Println("Timeouts", timeouts)
	fmt.Println("Calls for each serie", *calls)
	fmt.Println("Workerpool size", *poolSize)
	fmt.Println("Total series count", *runTimes)
	fmt.Println("Average time", timesAll/(float64)(len(times)), "ms")
}

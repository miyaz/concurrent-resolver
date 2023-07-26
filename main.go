package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	domain := "www.google.com"
	concurrency := 10
	count := 1000
	for i, v := range os.Args {
		if i == 1 {
			domain = v
		}
		if i == 2 {
			if num, err := strconv.Atoi(v); err != nil {
				fmt.Printf("%s %v", v, err)
				os.Exit(1)
			} else {
				concurrency = num
			}
		}
		if i == 3 {
			if num, err := strconv.Atoi(v); err != nil {
				//fmt.Printf("%s %v", v, err)
				os.Exit(1)
			} else {
				count = num
			}
		}
	}
	fmt.Printf("target domain:%s concurrency:%d count:%d\n", domain, concurrency, count)

	elapsed := concurrentResolver(domain, concurrency, count)
	fmt.Printf("count: %d elapsed(ms): %d, qps:%.3f\n", count, elapsed, float64(count)/(float64(elapsed)/1000))
}

// Result ... result
type Result struct {
	Elapsed int64
	IPs     []string
}

func resolve(domain string) (*Result, error) {
	//	addr, err := net.ResolveIPAddr("ip", domain)
	startTime := time.Now().UTC()
	startTimeMS := startTime.UnixNano() / int64(time.Millisecond)
	addrs, err := net.LookupHost(domain)
	endTime := time.Now().UTC()
	endTimeMS := endTime.UnixNano() / int64(time.Millisecond)
	//fmt.Printf("%s %s (%d)\n", time2str(startTime), time2str(endTime), endTimeMS-startTimeMS)
	if err != nil {
		fmt.Println("Resolve error ", err)
		return nil, err
	}
	return &Result{Elapsed: (endTimeMS - startTimeMS) / 2, IPs: addrs}, nil
}

func time2str(t time.Time) string {
	//return t.Format("2006-01-02T15:04:05")
	return t.Format("2006-01-02T15:04:05.999")
}

func concurrentResolver(domain string, concurrency int, count int) int {
	limit := make(chan int, concurrency)
	var wg sync.WaitGroup
	startTime := time.Now().UTC()
	startTimeMS := startTime.UnixNano() / int64(time.Millisecond)
	for i := 0; i < count/2; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, limit chan int) {
			defer wg.Done()
			limit <- 1
			if result, err := resolve(domain); err == nil {
				//fmt.Printf("%v elapsed(ms):%d\n", result.IPs, result.Elapsed)
				fmt.Printf("%d\n", result.Elapsed)
			} else {
				fmt.Println(err)
			}
			<-limit
		}(&wg, limit)
	}
	wg.Wait()
	endTime := time.Now().UTC()
	endTimeMS := endTime.UnixNano() / int64(time.Millisecond)
	return int(endTimeMS - startTimeMS)
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type worker interface {
	work(chan interface{}, *sync.WaitGroup)
}

type URLChecker struct {
	timeout int
}

func main() {
	workers := flag.Int("workers", 5, "Number of workers to use.")
	file := flag.String("file", "", "The input file to use. If none is supplied, stdin is assumed.")
	timeout := flag.Int("timeout", 10, "Max number of seconds to wait for a response from the URL.")
	var worker worker
	flag.Parse()

	work := make(chan interface{})

	if *file == "" {
		go feedURLs(work, os.Stdin)
	} else {
		urls, err := os.Open(*file)
		if err != nil {
			log.Fatal("unable to open input file:", err)
		}
		defer urls.Close()
		go feedURLs(work, urls)
	}

	var wg sync.WaitGroup
	worker = URLChecker{timeout: *timeout}

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker.work(work, &wg)
	}

	wg.Wait()
}

func newClient(timeout int) *http.Client {
	return &http.Client{Timeout: time.Second * time.Duration(timeout)}
}

func (w URLChecker) work(urls chan interface{}, wg *sync.WaitGroup) {
	c := newClient(w.timeout)

	for url := range urls {
		start := time.Now()
		resp, err := c.Head(url.(string))
		if err != nil {
			log.Printf("unable to connect to %s: %s", url, err.Error())
			continue
		}

		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		fmt.Printf("%s %d %6.3fms\n", url, resp.StatusCode, time.Since(start).Seconds()*1000)
	}
	wg.Done()
}

func feedURLs(work chan interface{}, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		work <- scanner.Text()
	}
	close(work)
}

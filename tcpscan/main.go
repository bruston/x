package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	flagMaxConn := flag.Int("max_con", 5, "The maximum number of simultanious connection attempts to make.")
	flagHost := flag.String("host", "localhost", "The hostname/IP address to scan.")
	flagPortMin := flag.Uint("min", 1, "The scanning starting port.")
	flagPortMax := flag.Uint("max", 1024, "The scanning end port.")
	flagTimeout := flag.Int("timeout", 5, "Time in seconds to wait for a port probe before giving up.")
	flag.Parse()

	if *flagHost == "" {
		flag.Usage()
		return
	}
	fmt.Printf("Scanning ports %d-%d on %s...\n", *flagPortMin, *flagPortMax, *flagHost)
	ports := feedPorts(*flagPortMin, *flagPortMax)

	var wg sync.WaitGroup
	for i := 0; i < *flagMaxConn; i++ {
		wg.Add(1)
		go checker(*flagHost, ports, &wg, time.Duration(*flagTimeout)*time.Second)
	}

	wg.Wait()

	fmt.Println("Done.")
}

func feedPorts(start, end uint) chan uint {
	portCh := make(chan uint)

	go func() {
		for i := start; i < end; i++ {
			portCh <- i
		}
		close(portCh)
	}()

	return portCh
}

func checker(host string, ports chan uint, wg *sync.WaitGroup, timeout time.Duration) {
	for port := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(int(port))), timeout)
		if err != nil {
			continue
		}
		conn.Close()

		fmt.Printf("%d\n", port)
	}

	wg.Done()
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:2202", "The host and port to listen on.")
	flag.Parse()

	server, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := server.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	fmt.Fprintln(c, "hello!")
	c.Close()
}

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":54321")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	file, err := os.Open("data.txt")
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()
	n, err := io.Copy(conn, file)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("Wrote %d bytes\n", n)
}

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:54321")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	file, err := os.Create("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	n, err := io.Copy(file, conn)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("Wrote %d bytes\n", n)
}

package main

import (
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:43249")
	if err != nil {
		log.Fatal(err)
	}
	uconn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer uconn.Close()
	buf := make([]byte, 512)
	for {
		n, raddr, err := uconn.ReadFromUDP(buf)
		if err != nil {
			log.Print(err)
			continue
		}
		if _, err := uconn.WriteToUDP(buf[0:n], raddr); err != nil {
			log.Print(err)
		}
	}
}

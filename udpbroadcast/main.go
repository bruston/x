package main

import (
	"log"
	"net"
	"time"
)

func main() {
	const port = 43249
	const sleepTime = time.Second * 30
	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	if err != nil {
		log.Print(err)
		return
	}
	go listen(socket)
	broadcast([]byte{1}, port)
	for {
		time.Sleep(sleepTime)
		if _, err := broadcast([]byte{1}, port); err != nil {
			log.Print(err)
			continue
		}
	}
}

func listen(socket *net.UDPConn) {
	for {
		data := make([]byte, 1)
		_, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			return
		}
		log.Printf("%s: %s", remoteAddr, data)
	}
}

func broadcast(payload []byte, port int) (int, error) {
	broadcastIP := net.IPv4(255, 255, 255, 255)
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   broadcastIP,
		Port: port,
	})
	if err != nil {
		return 0, err
	}
	socket.SetWriteDeadline(time.Now().Add(time.Second * 2))
	n, err := socket.Write([]byte{1})
	socket.Close()
	return n, err
}

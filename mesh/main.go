package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/bruston/lru"
)

type message struct {
	ID   string
	Addr string
	Body []byte
}

type server struct {
	ip    string
	port  string
	p     peers
	cache lru.Cacher
}

func main() {
	peer := flag.String("peer", "", "peer host:port")
	localIP := flag.String("ip", "", "your public IP address")
	flag.Parse()
	server := &server{
		ip: *localIP,
		p: peers{
			m: make(map[string]chan<- message),
		},
		cache: lru.NewSafeCache(100, 5),
	}

	l, err := server.listen()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on", server.self())
	go server.dial(*peer)
	fmt.Println("Dialing peer:", *peer)
	go server.readInput()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Print(err)
			return
		}
		go server.serve(c)
	}
}

type peers struct {
	m  map[string]chan<- message
	mu sync.RWMutex
}

func (p *peers) add(addr string) <-chan message {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.m[addr]; ok {
		return nil
	}
	ch := make(chan message)
	p.m[addr] = ch
	log.Printf("added peer: %s", addr)
	return ch
}

func (p *peers) remove(addr string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	delete(p.m, addr)
	log.Printf("removed peer: %s", addr)
}

func (p *peers) list() []chan<- message {
	p.mu.RLock()
	defer p.mu.RUnlock()
	l := make([]chan<- message, 0, len(p.m))
	for _, ch := range p.m {
		l = append(l, ch)
	}
	return l
}

func (s *server) serve(c net.Conn) {
	defer c.Close()
	d := json.NewDecoder(c)
	m := message{}
	for {
		if err := d.Decode(&m); err != nil {
			log.Print(err)
			return
		}
		if s.cache.Exists(m.ID) {
			continue
		}
		s.cache.Put(m.ID, struct{}{})
		s.broadcast(m)
		fmt.Printf("%#v\n", m)
		go s.dial(m.Addr)
	}
}

func (s *server) dial(addr string) {
	if addr == s.self() || addr == "" {
		return
	}

	ch := s.p.add(addr)
	if ch == nil {
		return
	}
	defer s.p.remove(addr)

	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("error dialing peer %s - %s", addr, err)
		return
	}
	defer c.Close()

	e := json.NewEncoder(c)
	for m := range ch {
		if err := e.Encode(m); err != nil {
			log.Print(err)
			return
		}
	}
}

func (s *server) broadcast(m message) {
	for _, ch := range s.p.list() {
		select {
		case ch <- m:
		default:
			// drop message
		}
	}
}

func (s *server) listen() (net.Listener, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	_, p, _ := net.SplitHostPort(l.Addr().String())
	s.port = p
	return l, err
}

func randID() string {
	b := make([]byte, 16)
	n, _ := rand.Read(b)
	return fmt.Sprintf("%x", b[:n])
}

func (s *server) readInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		m := message{
			ID:   randID(),
			Addr: s.self(),
			Body: scanner.Bytes(),
		}
		s.broadcast(m)
	}
}

func (s *server) self() string {
	return s.ip + ":" + s.port
}

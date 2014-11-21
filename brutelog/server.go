package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"net"

	"code.google.com/p/go.crypto/ssh"
)

type passCallback func(ssh.ConnMetadata, []byte) (*ssh.Permissions, error)

func main() {
	key := flag.String("key", "id_rsa", "The private key file to use.")
	listen := flag.String("listen", ":22", "The host and port to listen on. Example: :22")
	flag.Parse()

	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal("failed to start server: %s", err.Error())
	}

	config := newServerCfg(*key, passwordCallback)
	for {
		nconn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(nconn, config)
	}
}

func newServerCfg(key string, callback passCallback) *ssh.ServerConfig {
	privateBytes, err := ioutil.ReadFile(key)
	if err != nil {
		log.Fatalf("unable to read from private key file: %s", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("unable to parse private key: %s", err)
	}

	config := &ssh.ServerConfig{PasswordCallback: callback}
	config.AddHostKey(private)
	return config
}

func passwordCallback(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	fmt.Printf("%s\n", pass)
	return nil, fmt.Errorf("password rejected for %q", c.User())
}

func handleConn(c net.Conn, cfg *ssh.ServerConfig) {
	defer c.Close()
	_, _, _, err := ssh.NewServerConn(c, cfg)
	if err != nil && err != io.EOF {
		log.Print(err)
		return
	}
}

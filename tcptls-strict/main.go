package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	flagListen := flag.String("listen", "127.0.0.1:2022", "The host and port to listen on.")
	flagCert := flag.String("cert", "cert.pem", "Certificate file.")
	flagKey := flag.String("key", "key.pem", "Key file.")
	flag.Parse()

	cert, err := tls.LoadX509KeyPair(*flagCert, *flagKey)
	fatalOnError(err)

	cfg := newRestrictedCfg(cert)
	listener, err := tls.Listen("tcp", *flagListen, &cfg)
	fatalOnError(err)

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(c)
	}
}

func fatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleConn(c net.Conn) {
	fmt.Fprintln(c, "hello!")
	c.Close()
}
func newRestrictedCfg(cert tls.Certificate) tls.Config {
	return tls.Config{
		Certificates:           []tls.Certificate{cert},
		MinVersion:             tls.VersionTLS12,
		SessionTicketsDisabled: true,
		CipherSuites: []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

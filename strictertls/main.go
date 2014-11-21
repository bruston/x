package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

func main() {
	cert := flag.String("cert", "", "A PEM encoded certificate file")
	key := flag.String("key", "", "A PEM encoded private key file")
	listen := flag.String("listen", ":8090", "The host and port to listen on")
	flag.Parse()

	// Support only TLS1.0, TLS1.1 and TLS1.2
	cfg := &tls.Config{MinVersion: tls.VersionTLS10}
	server := &http.Server{Addr: *listen, TLSConfig: cfg, Handler: http.DefaultServeMux}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})

	log.Printf("starting webserver on %s\n", *listen)

	if err := server.ListenAndServeTLS(*cert, *key); err != nil {
		log.Fatal(err)
	}
}
